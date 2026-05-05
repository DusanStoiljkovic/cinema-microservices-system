package dto

type HallRequest struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Capacity uint   `json:"capacity"`
}
