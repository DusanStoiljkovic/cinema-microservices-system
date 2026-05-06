package models

import "time"

type Projection struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	MovieID   uint      `json:"movie_id" gorm:"not null;index"`
	HallID    uint      `json:"hall_id" gorm:"not null;index"`
	StartTime time.Time `json:"start_time" gorm:"not null"`
	EndTime   time.Time `json:"end_time" gorm:"not null"`
	Price     float64   `json:"price" gorm:"not null"`

	Tickets []Ticket `json:"tickets" gorm:"foreignKey:ProjectionID"`

	Hall Hall `json:"hall" gorm:"foreignKey:HallID"`
}
