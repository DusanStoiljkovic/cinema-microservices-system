package dto

type Status string

const (
	Reserved  Status = "reserved"
	Paid      Status = "paid"
	Cancelled        = "cancelled"
)

type OrderRequest struct {
	Tickets []TicketRequest `json:"tickets"`
}
