package models

// SysConfig is a key-value parameter (参考 go-admin sys_config 极简版).
// ConfigType: Y=内置, N=自定义; IsFrontend: 1=前台可读, 0=仅后台.
type SysConfig struct {
	BaseModel
	ControlBy
	ConfigName  string `gorm:"size:128" json:"configName"`
	ConfigKey   string `gorm:"size:128;uniqueIndex;not null" json:"configKey"`
	ConfigValue string `gorm:"size:1024" json:"configValue"`
	ConfigType  string `gorm:"size:8;default:N" json:"configType"`
	IsFrontend  int    `gorm:"default:0" json:"isFrontend"`
	Remark      string `gorm:"size:255" json:"remark"`
}

func (SysConfig) TableName() string { return "sys_config" }
