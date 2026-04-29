package models

import "time"

type Genre struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	Movies []Movie `gorm:"many2many:movie_genres;"`
}
