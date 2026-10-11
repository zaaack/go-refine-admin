package models

// Upload is a stored file record.
type Upload struct {
	BaseModel
	ControlBy
	Name string `gorm:"size:255" json:"name"`
	Url  string `gorm:"size:512" json:"url"`
	MD5  string `gorm:"size:64;index" json:"md5"`
	Size int64  `json:"size"`
	Ext  string `gorm:"size:16" json:"ext"`
}

func (Upload) TableName() string { return "sys_upload" }
