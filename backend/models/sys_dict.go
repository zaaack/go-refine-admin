package models

// SysDict is a dictionary (type) grouping several details.
type SysDict struct {
	BaseModel
	ControlBy
	Name   string `gorm:"size:64;not null" json:"name"`
	Type   string `gorm:"size:64;uniqueIndex;not null" json:"type"` // dict type key
	Status int    `gorm:"default:1;index" json:"status"`
	Remark string `gorm:"size:255" json:"remark"`
}

func (SysDict) TableName() string { return "sys_dict" }

// SysDictDetail is one key/value entry of a dictionary.
type SysDictDetail struct {
	BaseModel
	ControlBy
	DictID uint   `gorm:"index;not null" json:"dictId"`
	Label  string `gorm:"size:64" json:"label"`
	Value  string `gorm:"size:64" json:"value"`
	Sort   int    `gorm:"default:0" json:"sort"`
	Status int    `gorm:"default:1;index" json:"status"`
	Remarks string `gorm:"size:255" json:"remarks"`
}

func (SysDictDetail) TableName() string { return "sys_dict_detail" }
