package handler

import (
	"context"
	"encoding/json"
	"movie-service/internal/models"
	"net/http"
)

type MovieService interface {
	GetMovies(ctx context.Context,
		limit, offset int,
		sort string,
		genre string,
		minYear, maxYear int,
		minRating float64) ([]*models.Movie, error)
	GetMovieByID(ctx context.Context, id uint) (*models.Movie, error)
	CreateMovie(ctx context.Context, movie *models.Movie) (*models.Movie, error)
	UpdateMovie(ctx context.Context, id uint, movie *models.Movie) (*models.Movie, error)
	DeleteMovie(ctx context.Context, id uint) error
}

type MovieHandler struct {
	service MovieService
}

func NewMovieHandler(service MovieService) *MovieHandler {
	return &MovieHandler{service: service}
}

func (handler *MovieHandler) HandleGetMovies(w http.ResponseWriter, r *http.Request) error {
	movies, err := handler.service.GetMovies(r.Context(), 0, 0, "", "", 2000, 2024, 0.0)
	if err != nil {
		http.Error(w, "Failed to get movies", http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(&movies)
	return nil
}

func (handler *MovieHandler) HandleGetMovieByID(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (handler *MovieHandler) HandleCreateMovie(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (handler *MovieHandler) HandleUpdateMovie(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (handler *MovieHandler) HandleDeleteMovie(w http.ResponseWriter, r *http.Request) error {
	return nil
}
