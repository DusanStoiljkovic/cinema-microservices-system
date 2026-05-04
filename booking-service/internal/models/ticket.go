package models

type TicketStatus string

const (
	TicketStatusReserved  TicketStatus = "reserved"
	TicketStatusPaid      TicketStatus = "paid"
	TicketStatusCancelled TicketStatus = "cancelled"
)

type Ticket struct {
	ID           uint         `json:"id" gorm:"primaryKey"`
	UserID       uint         `json:"user_id" gorm:"not null;index"`
	ProjectionID uint         `json:"projection_id" gorm:"not null;uniqueIndex:idx_projection_seat"`
	SeatNumber   uint         `json:"seat_number" gorm:"not null;uniqueIndex:idx_projection_seat"`
	Status       TicketStatus `json:"status" gorm:"type:varchar(20);not null;default:'reserved'"`

	Projection Projection `json:"projection" gorm:"foreignKey:ProjectionID"`
}
