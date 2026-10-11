package models

// SysMenu is a navigation/menu node. Type: M=directory, C=menu(page), F=button.
type SysMenu struct {
	BaseModel
	ControlBy
	Name      string `gorm:"size:64;not null" json:"name"`    // route name / component key
	Title     string `gorm:"size:64;not null" json:"title"`   // i18n title
	Icon      string `gorm:"size:64" json:"icon"`
	Path      string `gorm:"size:128" json:"path"`
	Component string `gorm:"size:128" json:"component"`
	ParentID  uint   `gorm:"default:0;index" json:"parentId"`
	Sort      int    `gorm:"default:0" json:"sort"`
	Type      string `gorm:"size:1;default:M" json:"type"` // M | C | F
	Permission string `gorm:"size:128" json:"permission"`
	Status    int    `gorm:"default:1;index" json:"status"`
	Hidden    int    `gorm:"default:0" json:"hidden"`
}

func (SysMenu) TableName() string { return "sys_menu" }
