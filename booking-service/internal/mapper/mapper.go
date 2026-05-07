package mapper

import (
	"booking-service/internal/dto"
	"booking-service/internal/models"
)

func HallFromRequest(req *dto.HallRequest) *models.Hall {
	return &models.Hall{
		Name:     req.Name,
		Location: req.Location,
		Capacity: req.Capacity,
	}
}

func ProjectionFromRequest(req *dto.ProjectionRequest) *models.Projection {
	return &models.Projection{
		MovieID:   req.MovieID,
		HallID:    req.HallID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Price:     req.Price,
	}
}

func TicketFromRequest(req *dto.TicketRequest) *models.Ticket {
	return &models.Ticket{
		ProjectionID: req.ProjectionID,
		SeatNumber:   req.SeatNumber,
	}
}

func OrderFromRequest(req *dto.OrderRequest) *models.Order {
	tickets := make([]models.Ticket, 0, len(req.Tickets))

	for _, val := range req.Tickets {
		tickets = append(tickets, *TicketFromRequest(&val))
	}

	return &models.Order{
		Tickets: tickets,
	}
}
