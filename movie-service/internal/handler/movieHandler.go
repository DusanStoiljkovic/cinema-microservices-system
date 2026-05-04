package handler

import (
	"context"
	"encoding/json"
	"movie-service/internal/dto"
	"movie-service/internal/mapper"
	"movie-service/internal/models"
	"movie-service/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type MovieService interface {
	GetMovies(
		ctx context.Context,
		limit, offset int,
		sort string,
		genre string,
		minYear, maxYear int,
		minRating float64,
	) ([]*models.Movie, error)

	GetMovieByID(ctx context.Context, id uint) (*models.Movie, error)
	GetRelationsByMovieID(ctx context.Context, id uint) ([]models.Genre, error)
	CreateMovie(ctx context.Context, movie *models.Movie) (*models.Movie, error)
	CreateRelation(ctx context.Context, movieID, genreID uint) (*models.Movie, error)
	UpdateMovie(ctx context.Context, id uint, movie *models.Movie) (*models.Movie, error)
	DeleteMovie(ctx context.Context, id uint) error
	DeleteRelation(ctx context.Context, movieID, genreID uint) error
}

type MovieHandler struct {
	service MovieService
}

func NewMovieHandler(service MovieService) *MovieHandler {
	return &MovieHandler{service: service}
}

func (handler *MovieHandler) HandleGetMovies(w http.ResponseWriter, r *http.Request) error {
	limit := getIntQuery(r, "limit", 20)
	offset := getIntQuery(r, "offset", 0)
	sort := r.URL.Query().Get("sort")
	genre := r.URL.Query().Get("genre")
	minYear := getIntQuery(r, "min_year", 0)
	maxYear := getIntQuery(r, "max_year", 0)
	minRating := getFloatQuery(r, "min_rating", 0.0)

	movies, err := handler.service.GetMovies(
		r.Context(),
		limit,
		offset,
		sort,
		genre,
		minYear,
		maxYear,
		minRating,
	)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, movies)
}

func (handler *MovieHandler) HandleGetMovieByID(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIDParam(r, "id")
	if err != nil {
		return utils.ErrInvalidInput
	}

	movie, err := handler.service.GetMovieByID(r.Context(), id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, movie)
}

func (handler *MovieHandler) HandleGetRelationsByMovieID(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIDParam(r, "id")
	if err != nil {
		return utils.ErrInvalidInput
	}

	relations, err := handler.service.GetRelationsByMovieID(r.Context(), id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, relations)
}

func (handler *MovieHandler) HandleCreateMovie(w http.ResponseWriter, r *http.Request) error {
	var req dto.MovieRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.ErrInvalidInput
	}

	movie := mapper.MovieFromRequest(req)

	createdMovie, err := handler.service.CreateMovie(r.Context(), movie)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusCreated, createdMovie)
}

func (handler *MovieHandler) HandleCreateRelation(w http.ResponseWriter, r *http.Request) error {
	movieId, err := parseIDParam(r, "movieId")
	if err != nil {
		return utils.ErrInvalidInput
	}

	genreId, err := parseIDParam(r, "genreId")
	if err != nil {
		return utils.ErrInvalidInput
	}

	movie, err := handler.service.CreateRelation(r.Context(), movieId, genreId)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusCreated, movie)
}

func (handler *MovieHandler) HandleUpdateMovie(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIDParam(r, "id")
	if err != nil {
		return utils.ErrInvalidInput
	}

	var req dto.MovieRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return utils.ErrInvalidInput
	}

	movie := mapper.MovieFromRequest(req)

	updatedMovie, err := handler.service.UpdateMovie(r.Context(), id, movie)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, updatedMovie)
}

func (handler *MovieHandler) HandleDeleteMovie(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIDParam(r, "id")
	if err != nil {
		return utils.ErrInvalidInput
	}

	if err := handler.service.DeleteMovie(r.Context(), id); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (handler *MovieHandler) HandleDeleteRelation(w http.ResponseWriter, r *http.Request) error {
	movieId, err := parseIDParam(r, "movieId")
	if err != nil {
		return utils.ErrInvalidInput
	}

	genreId, err := parseIDParam(r, "genreId")
	if err != nil {
		return utils.ErrInvalidInput
	}

	err = handler.service.DeleteRelation(r.Context(), movieId, genreId)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusNoContent, "")
}

func parseIDParam(r *http.Request, param string) (uint, error) {
	value := chi.URLParam(r, param)

	id, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func getIntQuery(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return parsedValue
}

func getFloatQuery(r *http.Request, key string, defaultValue float64) float64 {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}

	parsedValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}

	return parsedValue
}
