package models

// SysRole is an RBAC role. MenuIDs is a comma-separated list of granted menu ids.
type SysRole struct {
	BaseModel
	ControlBy
	Name     string `gorm:"size:64;not null" json:"name"`
	RoleKey  string `gorm:"size:64;uniqueIndex;not null" json:"roleKey"`
	Status   int    `gorm:"default:1;index" json:"status"`
	Sort     int    `gorm:"default:0" json:"sort"`
	DataScope int   `gorm:"default:1" json:"dataScope"` // 1 all, 2 custom, 3 self
	Remark   string `gorm:"size:255" json:"remark"`
	MenuIDs  string `gorm:"type:text" json:"menuIds"`
}

func (SysRole) TableName() string { return "sys_role" }
