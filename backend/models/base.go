package models

import "time"

// Entity is the interface every actions-mode model must satisfy.
type Entity interface {
	TableName() string
	GetId() interface{}
}

// BaseModel provides the standard primary key + timestamp fields.
// Embed it in your models to get CRUD + audit fields for free.
type BaseModel struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// GetId implements Entity.
func (b BaseModel) GetId() interface{} { return b.ID }

// ControlBy records which user created/updated a row (data audit).
type ControlBy struct {
	CreateBy uint `gorm:"index" json:"createBy,omitempty"`
	UpdateBy uint `gorm:"index" json:"updateBy,omitempty"`
}

func (c *ControlBy) SetCreateBy(uid uint) { c.CreateBy = uid }
func (c *ControlBy) SetUpdateBy(uid uint) { c.UpdateBy = uid }
