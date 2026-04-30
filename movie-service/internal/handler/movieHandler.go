package handler

import (
	"context"
	"encoding/json"
	"errors"
	"movie-service/internal/models"
	servicepkg "movie-service/internal/service"
	"movie-service/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	id, err := parseIDParam(r, "id")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid movie id")
		return err
	}

	movie, err := handler.service.GetMovieByID(r.Context(), id)
	if err != nil {
		return HandlerError(w, err)
	}

	return utils.WriteJSON(w, http.StatusOK, movie)
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

func parseIDParam(r *http.Request, param string) (uint, error) {
	value := chi.URLParam(r, param)

	id, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func HandlerError(w http.ResponseWriter, err error) error {
	switch {
	case errors.Is(err, servicepkg.ErrInvalidInput):
		utils.WriteError(w, http.StatusBadRequest, "invalid input")
	case errors.Is(err, servicepkg.ErrNotFound):
		utils.WriteError(w, http.StatusNotFound, "resource not found")
	case errors.Is(err, servicepkg.ErrConflict):
		utils.WriteError(w, http.StatusConflict, "resource conflict")
	default:
		utils.WriteError(w, http.StatusInternalServerError, "internal server error")
	}

	return err
}
