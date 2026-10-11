package db

import (
	"log"
	"strconv"
	"time"

	"go-refine-admin/backend/models"
	"gorm.io/gorm"
)

// Retention config keys and floor. Logs are never kept less than minRetentionDays.
const (
	operRetentionKey  = "sys.log.oper-retention-days"
	loginRetentionKey = "sys.log.login-retention-days"
	minRetentionDays  = 30
)

// ensureRetentionConfigs seeds the retention parameters when missing (fresh and
// existing databases alike). Values are editable in 参数管理 (sys/config).
func ensureRetentionConfigs(gdb *gorm.DB) error {
	defs := []models.SysConfig{
		{ConfigName: "操作日志保留天数", ConfigKey: operRetentionKey, ConfigValue: "30", ConfigType: "Y", IsFrontend: 0, Remark: "低于30按30天执行"},
		{ConfigName: "登录日志保留天数", ConfigKey: loginRetentionKey, ConfigValue: "30", ConfigType: "Y", IsFrontend: 0, Remark: "低于30按30天执行"},
	}
	for i := range defs {
		var c models.SysConfig
		if err := gdb.Where("config_key = ?", defs[i].ConfigKey).First(&c).Error; err != nil {
			if err := gdb.Create(&defs[i]).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// retentionDays reads a retention parameter, falling back to the floor on any
// missing/invalid/too-small value.
func retentionDays(gdb *gorm.DB, key string) int {
	var c models.SysConfig
	if err := gdb.Where("config_key = ?", key).First(&c).Error; err != nil {
		return minRetentionDays
	}
	n, err := strconv.Atoi(c.ConfigValue)
	if err != nil || n < minRetentionDays {
		return minRetentionDays
	}
	return n
}

// purgeLogs deletes audit/login rows older than their retention window.
func purgeLogs(gdb *gorm.DB) {
	cut := func(days int) time.Time { return time.Now().AddDate(0, 0, -days) }
	if err := gdb.Where("created_at < ?", cut(retentionDays(gdb, operRetentionKey))).
		Delete(&models.SysOperLog{}).Error; err != nil {
		log.Printf("purge oper-log: %v", err)
	}
	if err := gdb.Where("created_at < ?", cut(retentionDays(gdb, loginRetentionKey))).
		Delete(&models.SysLoginLog{}).Error; err != nil {
		log.Printf("purge login-log: %v", err)
	}
}

// StartLogRetentionJanitor purges expired logs at boot and every 24h.
// It runs in its own goroutine; callers must invoke it with `go`.
func StartLogRetentionJanitor(gdb *gorm.DB) {
	purgeLogs(gdb)
	for range time.Tick(24 * time.Hour) {
		purgeLogs(gdb)
	}
}
