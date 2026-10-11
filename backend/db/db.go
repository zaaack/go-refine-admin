package db

import (
	"fmt"

	"enterprise-admin/backend/config"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open builds a *gorm.DB from configuration. Supports sqlite (default, zero-config),
// mysql and postgres. The caller is responsible for running migrations/seed.
func Open(cfg config.DatabaseConfig) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch cfg.Driver {
	case "mysql":
		dialector = mysql.Open(cfg.MySQL.DSN())
	case "postgres":
		dialector = postgres.Open(cfg.Postgres.DSN())
	case "sqlite", "":
		dialector = sqlite.Open(cfg.SQLite.Path)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	gdb, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	return gdb, nil
}
