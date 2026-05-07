package mapper

import (
	"booking-service/internal/dto"
	"booking-service/internal/models"
)

func OrderToResponse(order *models.Order) *dto.OrderResponse {
	if order == nil {
		return nil
	}

	var tickets []*models.Ticket
	for _, val := range order.Tickets {
		tickets = append(tickets, &val)
	}

	return &dto.OrderResponse{
		ID:         order.ID,
		UserID:     order.UserID,
		Status:     order.Status,
		TotalPrice: order.TotalPrice,
		Tickets:    TicketsToResponse(tickets),
	}
}

func OrdersToResponse(orders []*models.Order) []dto.OrderResponse {
	responses := make([]dto.OrderResponse, 0, len(orders))

	for _, order := range orders {
		response := OrderToResponse(order)
		if response != nil {
			responses = append(responses, *response)
		}
	}

	return responses
}

func TicketToResponse(ticket *models.Ticket) *dto.TicketResponse {
	if ticket == nil {
		return nil
	}

	response := &dto.TicketResponse{
		ID:           ticket.ID,
		ProjectionID: ticket.ProjectionID,
		OrderID:      ticket.OrderID,
		SeatNumber:   ticket.SeatNumber,
		CreatedAt:    ticket.CreatedAt,
	}

	if ticket.Projection.ID != 0 {
		response.Projection = ProjectionToResponse(&ticket.Projection)
	}

	return response
}

func TicketsToResponse(tickets []*models.Ticket) []dto.TicketResponse {
	responses := make([]dto.TicketResponse, 0, len(tickets))

	for _, ticket := range tickets {
		response := TicketToResponse(ticket)
		if response != nil {
			responses = append(responses, *response)
		}
	}

	return responses
}

func ProjectionToResponse(projection *models.Projection) *dto.ProjectionResponse {
	if projection == nil {
		return nil
	}

	response := &dto.ProjectionResponse{
		ID:        projection.ID,
		MovieID:   projection.MovieID,
		HallID:    projection.HallID,
		StartTime: projection.StartTime,
		EndTime:   projection.EndTime,
		Price:     projection.Price,
	}

	if projection.Hall.ID != 0 {
		response.Hall = HallToResponse(&projection.Hall)
	}

	return response
}

func ProjectionsToResponse(projections []*models.Projection) []dto.ProjectionResponse {
	responses := make([]dto.ProjectionResponse, 0, len(projections))

	for _, projection := range projections {
		response := ProjectionToResponse(projection)
		if response != nil {
			responses = append(responses, *response)
		}
	}

	return responses
}

func HallToResponse(hall *models.Hall) *dto.HallResponse {
	if hall == nil {
		return nil
	}

	return &dto.HallResponse{
		ID:       hall.ID,
		Name:     hall.Name,
		Location: hall.Location,
		Capacity: hall.Capacity,
	}
}

func HallsToResponse(halls []*models.Hall) []dto.HallResponse {
	responses := make([]dto.HallResponse, 0, len(halls))

	for _, hall := range halls {
		response := HallToResponse(hall)
		if response != nil {
			responses = append(responses, *response)
		}
	}

	return responses
}
