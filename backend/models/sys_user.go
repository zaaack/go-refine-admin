package models

import "golang.org/x/crypto/bcrypt"

// SysUser is the system user (login account).
type SysUser struct {
	BaseModel
	ControlBy
	Username string `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:128;not null" json:"-"`
	// PlainPassword is bind-only: accepted from JSON on create/update, hashed
	// into Password by a hook, never persisted or serialized back.
	PlainPassword string `gorm:"-" json:"password,omitempty"`
	NickName      string `gorm:"size:64" json:"nickName"`
	Phone         string `gorm:"size:32" json:"phone"`
	Email         string `gorm:"size:64" json:"email"`
	Avatar        string `gorm:"size:255" json:"avatar"`
	Sex           int    `gorm:"default:0" json:"sex"`
	Status        int    `gorm:"default:1;index" json:"status"` // 1 enabled, 0 disabled
	DeptID        uint   `gorm:"index" json:"deptId"`
	RoleID        uint   `gorm:"index" json:"roleId"`
	Remark        string `gorm:"size:255" json:"remark"`
}

func (SysUser) TableName() string { return "sys_user" }

// SetPassword hashes and stores the password.
func (u *SysUser) SetPassword(plain string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// ComparePassword validates a plaintext password.
func (u *SysUser) ComparePassword(plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain)) == nil
}
