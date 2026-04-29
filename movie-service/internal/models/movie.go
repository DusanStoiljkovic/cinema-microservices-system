package models

import "time"

type Movie struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"not null"`
	Description string    `gorm:"not null"`
	Year        int       `gorm:"not null"`
	ImageUrl    string    `gorm:"not null"`
	Duration    int       `gorm:"not null"`
	Rating      float64   `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`

	Genres []Genre `gorm:"many2many:movie_genres;"`
}
