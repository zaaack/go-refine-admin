package models

import "time"

// SysUserToken is the server-side session record for a JWT (enables force-logout).
// A JWT is only valid when its JTI exists here with Status == 0 and ExpireAt in future.
type SysUserToken struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"index" json:"userId"`
	Username  string    `gorm:"size:64;index" json:"username"`
	JTI       string    `gorm:"size:64;uniqueIndex" json:"jti"`
	IP        string    `gorm:"size:64" json:"ip"`
	Location  string    `gorm:"size:128" json:"location"`
	Browser   string    `gorm:"size:64" json:"browser"`
	OS        string    `gorm:"size:64" json:"os"`
	Device    string    `gorm:"size:64" json:"device"`
	UA        string    `gorm:"type:text" json:"ua"`
	LoginTime time.Time `json:"loginTime"`
	ExpireAt  time.Time `gorm:"index" json:"expireAt"`
	Status    int       `gorm:"default:0;index" json:"status"` // 0 online, 1 offline
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (SysUserToken) TableName() string { return "sys_user_token" }

// GetId implements Entity.
func (m SysUserToken) GetId() interface{} { return m.ID }
