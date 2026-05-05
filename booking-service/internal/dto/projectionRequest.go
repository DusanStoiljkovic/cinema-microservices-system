package dto

import "time"

type ProjectionRequest struct {
	MovieID   uint      `json:"movie_id"`
	HallID    uint      `json:"hall_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Price     float64   `json:"price"`
}
