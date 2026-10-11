package models

import "time"

// SysOperLog records operation audit entries (login, CRUD, errors).
type SysOperLog struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title      string    `gorm:"size:255" json:"title"`
	Method     string    `gorm:"size:16" json:"method"`
	Path       string    `gorm:"size:255" json:"path"`
	IP         string    `gorm:"size:64" json:"ip"`
	UserID     uint      `gorm:"index" json:"userId"`
	Username   string    `gorm:"size:64" json:"username"`
	Status      int      `gorm:"default:0" json:"status"` // 0 success, 1 fail
	Duration   int64     `json:"duration"`                // ms
	ErrMsg     string    `gorm:"type:text" json:"errMsg"`
	Body       string    `gorm:"type:text" json:"body"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (SysOperLog) TableName() string { return "sys_oper_log" }
