package models

import "time"

// SysLoginLog records every login attempt (success/fail) with ip/device/location.
type SysLoginLog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"size:64;index" json:"username"`
	UserID    uint      `gorm:"index" json:"userId"`
	IP        string    `gorm:"size:64" json:"ip"`
	Location  string    `gorm:"size:128" json:"location"`
	Browser   string    `gorm:"size:64" json:"browser"`
	OS        string    `gorm:"size:64" json:"os"`
	Device    string    `gorm:"size:64" json:"device"`
	UA        string    `gorm:"type:text" json:"ua"`
	Status    int       `gorm:"default:0;index" json:"status"` // 0 success, 1 fail
	Msg       string    `gorm:"size:255" json:"msg"`
	CreatedAt time.Time `json:"createdAt"`
}

func (SysLoginLog) TableName() string { return "sys_login_log" }

// GetId implements Entity (allows Actions-mode reuse if needed).
func (m SysLoginLog) GetId() interface{} { return m.ID }
