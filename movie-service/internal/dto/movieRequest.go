package dto

type MovieRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Year        int     `json:"year"`
	ImageURL    string  `json:"image_url"`
	Duration    int     `json:"duration"`
	Rating      float64 `json:"rating"`
	GenreIDs    []uint  `json:"genre_ids"`
}


