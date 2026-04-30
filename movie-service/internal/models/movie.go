package models

import "time"

type Movie struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	Year        int       `json:"year" gorm:"not null"`
	ImageURL    string    `json:"image_url" gorm:"not null"`
	Duration    int       `json:"duration" gorm:"not null"`
	Rating      float64   `json:"rating" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Genres []Genre `json:"genres" gorm:"many2many:movie_genres;"`
}
