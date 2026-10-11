package models

// Product is the demo resource registered in Actions mode (zero-code CRUD).
type Product struct {
	BaseModel
	ControlBy
	Name     string  `gorm:"size:128;not null;index" json:"name"`
	SKU      string  `gorm:"size:64;uniqueIndex" json:"sku"`
	Category string  `gorm:"size:64;index" json:"category"`
	Price    float64 `gorm:"type:decimal(12,2);default:0" json:"price"`
	Stock    int     `gorm:"default:0" json:"stock"`
	Status   int     `gorm:"default:1;index" json:"status"` // 1 on sale, 0 off
	Remark   string  `gorm:"size:255" json:"remark"`
}

func (Product) TableName() string { return "demo_product" }
