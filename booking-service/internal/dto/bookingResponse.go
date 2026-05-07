package dto

import (
	"booking-service/internal/models"
	"time"
)

type HallResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Capacity uint   `json:"capacity"`
}

type ProjectionResponse struct {
	ID        uint      `json:"id"`
	MovieID   uint      `json:"movie_id"`
	HallID    uint      `json:"hall_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Price     float64   `json:"price"`

	Hall *HallResponse `json:"hall,omitempty"`
}

type TicketResponse struct {
	ID           uint      `json:"id"`
	ProjectionID uint      `json:"projection_id"`
	OrderID      uint      `json:"order_id"`
	SeatNumber   uint      `json:"seat_number"`
	CreatedAt    time.Time `json:"created_at"`

	Projection *ProjectionResponse `json:"projection,omitempty"`
}

type OrderResponse struct {
	ID         uint               `json:"id"`
	UserID     uint               `json:"user_id"`
	Status     models.OrderStatus `json:"status"`
	TotalPrice uint               `json:"total_price"`

	Tickets []TicketResponse `json:"tickets"`
}
