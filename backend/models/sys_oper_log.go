package models

import "time"

// SysOperLog records operation audit entries (write operations only: POST/PUT/DELETE).
type SysOperLog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string    `gorm:"size:255" json:"title"`
	Method    string    `gorm:"size:16" json:"method"`
	Path      string    `gorm:"size:255" json:"path"`
	IP        string    `gorm:"size:64" json:"ip"`
	Location  string    `gorm:"size:128" json:"location"`
	UserAgent string    `gorm:"size:512" json:"userAgent"`
	UserID    uint      `gorm:"index" json:"userId"`
	Username  string    `gorm:"size:64" json:"username"`
	Status    int       `gorm:"default:0" json:"status"` // 0 success, 1 fail
	StatusCode int      `json:"statusCode"`
	Duration  int64     `json:"duration"` // ms
	ErrMsg    string    `gorm:"type:text" json:"errMsg"`
	Body      string    `gorm:"type:text" json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}

func (SysOperLog) TableName() string { return "sys_oper_log" }

// GetId implements Entity.
func (m SysOperLog) GetId() interface{} { return m.ID }
