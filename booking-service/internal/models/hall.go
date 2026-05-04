package models

type Hall struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null;unique"`
	Location string `json:"location"`
	Capacity int    `json:"capacity" gorm:"not null"`
}
