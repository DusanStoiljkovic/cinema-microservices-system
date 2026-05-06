package dto

type TicketRequest struct {
	ProjectionID uint `json:"projection_id"`
	SeatNumber   uint `json:"seat_number"`
	OrderID      uint `json:"order_id"`
}
