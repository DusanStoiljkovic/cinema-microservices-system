package models

import "time"

type Ticket struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	ProjectionID uint      `json:"projection_id" gorm:"not null;uniqueIndex:idx_projection_seat"`
	OrderID      uint      `json:"order_id" gorm:"not null;"`
	SeatNumber   uint      `json:"seat_number" gorm:"not null;uniqueIndex:idx_projection_seat"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null"`

	Projection Projection `json:"projection" gorm:"foreignKey:ProjectionID"`
	Order      Order      `json:"order" gorm:"foreignKey:OrderID"`
}
