package models

// SysDept is an organization department tree node.
type SysDept struct {
	BaseModel
	ControlBy
	Name   string `gorm:"size:64;not null" json:"name"`
	ParentID uint `gorm:"default:0;index" json:"parentId"`
	Sort   int    `gorm:"default:0" json:"sort"`
	Leader string `gorm:"size:64" json:"leader"`
	Phone  string `gorm:"size:32" json:"phone"`
	Status int    `gorm:"default:1;index" json:"status"`
}

func (SysDept) TableName() string { return "sys_dept" }
