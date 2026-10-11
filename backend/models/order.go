package models

// Order is the demo resource implemented in regular mode, showcasing the
// type-safe query builder for complex business logic.
type Order struct {
	BaseModel
	ControlBy
	OrderNo string  `gorm:"size:64;uniqueIndex;not null" json:"orderNo"`
	UserID  uint    `gorm:"index" json:"userId"`
	Amount  float64 `gorm:"type:decimal(12,2);default:0" json:"amount"`
	Status  string  `gorm:"size:16;index" json:"status"` // pending|paid|shipped|done|cancel
	Channel string  `gorm:"size:32;index" json:"channel"`
	Remark  string  `gorm:"size:255" json:"remark"`
}

func (Order) TableName() string { return "demo_order" }

// OrderStatus constants (also used by the type-safe query builder).
const (
	OrderPending  = "pending"
	OrderPaid     = "paid"
	OrderShipped  = "shipped"
	OrderDone     = "done"
	OrderCanceled = "cancel"
)
