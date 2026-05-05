package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"movie-service/internal/dto"
	"movie-service/internal/mapper"
	"movie-service/internal/models"
	"movie-service/internal/utils"

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
		return utils.NewInvalidInput("Invalid movie id", err)
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
		return utils.NewInvalidInput("Invalid movie id", err)
	}

	relations, err := handler.service.GetRelationsByMovieID(r.Context(), id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, relations)
}

func (handler *MovieHandler) HandleCreateMovie(w http.ResponseWriter, r *http.Request) error {
	req := &dto.MovieRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.NewInvalidInput("Invalid request body", err)
	}

	movie := mapper.MovieFromRequest(*req)

	createdMovie, err := handler.service.CreateMovie(r.Context(), movie)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusCreated, createdMovie)
}

func (handler *MovieHandler) HandleCreateRelation(w http.ResponseWriter, r *http.Request) error {
	movieID, err := parseIDParam(r, "movieId")
	if err != nil {
		return utils.NewInvalidInput("Invalid movie id", err)
	}

	genreID, err := parseIDParam(r, "genreId")
	if err != nil {
		return utils.NewInvalidInput("Invalid genre id", err)
	}

	movie, err := handler.service.CreateRelation(r.Context(), movieID, genreID)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusCreated, movie)
}

func (handler *MovieHandler) HandleUpdateMovie(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIDParam(r, "id")
	if err != nil {
		return utils.NewInvalidInput("Invalid movie id", err)
	}

	req := &dto.MovieRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.NewInvalidInput("Invalid request body", err)
	}

	movie := mapper.MovieFromRequest(*req)

	updatedMovie, err := handler.service.UpdateMovie(r.Context(), id, movie)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, updatedMovie)
}

func (handler *MovieHandler) HandleDeleteMovie(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIDParam(r, "id")
	if err != nil {
		return utils.NewInvalidInput("Invalid movie id", err)
	}

	if err := handler.service.DeleteMovie(r.Context(), id); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (handler *MovieHandler) HandleDeleteRelation(w http.ResponseWriter, r *http.Request) error {
	movieID, err := parseIDParam(r, "movieId")
	if err != nil {
		return utils.NewInvalidInput("Invalid movie id", err)
	}

	genreID, err := parseIDParam(r, "genreId")
	if err != nil {
		return utils.NewInvalidInput("Invalid genre id", err)
	}

	if err := handler.service.DeleteRelation(r.Context(), movieID, genreID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func parseIDParam(r *http.Request, param string) (uint, error) {
	value := chi.URLParam(r, param)

	id, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("id must be greater than zero")
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
