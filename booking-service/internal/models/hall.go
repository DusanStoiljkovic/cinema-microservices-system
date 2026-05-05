package models

type Hall struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null;unique"`
	Location string `json:"location"`
	Capacity uint   `json:"capacity" gorm:"not null"`
}
