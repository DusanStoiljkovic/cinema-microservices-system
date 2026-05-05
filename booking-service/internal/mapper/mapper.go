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
