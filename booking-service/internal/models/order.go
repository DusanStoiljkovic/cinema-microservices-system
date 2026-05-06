package models

type OrderStatus string

const (
	OrderStatusReserved  OrderStatus = "reserved"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID         uint        `json:"id" gorm:"primaryKey"`
	UserID     uint        `json:"user_id" gorm:"not null;index"`
	Status     OrderStatus `json:"status" gorm:"type:enum('reserved','paid','cancelled');not null;default:'reserved'"`
	TotalPrice uint        `json:"total_price" gorm:"not null"`

	Tickets []Ticket `json:"tickets" gorm:"foreignKey:OrderID"`
}
