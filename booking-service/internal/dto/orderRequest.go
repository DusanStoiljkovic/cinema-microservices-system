package dto

type Status string

const (
	Reserved  Status = "reserved"
	Paid      Status = "paid"
	Cancelled        = "cancelled"
)

type OrderRequest struct {
	TotalPrice uint `json:"total_price" gorm:"not null"`

	Tickets []TicketRequest
}
