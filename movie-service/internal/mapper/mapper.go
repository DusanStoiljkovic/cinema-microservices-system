package mapper

import (
	"movie-service/internal/dto"
	"movie-service/internal/models"
)

func MovieFromRequest(req dto.MovieRequest) *models.Movie {
	movie := &models.Movie{
		Title:       req.Title,
		Description: req.Description,
		Year:        req.Year,
		ImageURL:    req.ImageURL,
		Duration:    req.Duration,
		Rating:      req.Rating,
	}

	if req.GenreIDs != nil {
		movie.Genres = make([]models.Genre, 0, len(req.GenreIDs))

		for _, genreID := range req.GenreIDs {
			movie.Genres = append(movie.Genres, models.Genre{ID: genreID})
		}
	}

	return movie
}
