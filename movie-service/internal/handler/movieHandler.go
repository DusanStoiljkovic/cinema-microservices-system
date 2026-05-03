package handler

import (
	"context"
	"encoding/json"
	"errors"
<<<<<<< HEAD
<<<<<<< HEAD
=======
	"movie-service/internal/dto"
	"movie-service/internal/mapper"
>>>>>>> feature/movieService
=======
	"movie-service/internal/dto"
	"movie-service/internal/mapper"
>>>>>>> da5f31b (feat(movie-service): implement genre management with repository, service, and handler layers; enhance movie handler and routes)
	"movie-service/internal/models"
	servicepkg "movie-service/internal/service"
	"movie-service/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type MovieService interface {
<<<<<<< HEAD
<<<<<<< HEAD
	GetMovies(ctx context.Context,
=======
	GetMovies(
		ctx context.Context,
>>>>>>> feature/movieService
=======
	GetMovies(
		ctx context.Context,
>>>>>>> da5f31b (feat(movie-service): implement genre management with repository, service, and handler layers; enhance movie handler and routes)
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
<<<<<<< HEAD
<<<<<<< HEAD
	movies, err := handler.service.GetMovies(r.Context(), 0, 0, "", "", 2000, 2024, 0.0)
=======
=======
>>>>>>> da5f31b (feat(movie-service): implement genre management with repository, service, and handler layers; enhance movie handler and routes)
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
<<<<<<< HEAD
>>>>>>> feature/movieService
=======
>>>>>>> da5f31b (feat(movie-service): implement genre management with repository, service, and handler layers; enhance movie handler and routes)
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
<<<<<<< HEAD
<<<<<<< HEAD
	return nil
=======
	var req dto.MovieRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	movie := mapper.MovieFromRequest(req)

	createdMovie, err := handler.service.CreateMovie(r.Context(), movie)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusCreated, createdMovie)
>>>>>>> da5f31b (feat(movie-service): implement genre management with repository, service, and handler layers; enhance movie handler and routes)
}

func (handler *MovieHandler) HandleUpdateMovie(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIDParam(r, "id")
	if err != nil {
		return err
	}

	var req dto.MovieRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	movie := mapper.MovieFromRequest(req)

	updatedMovie, err := handler.service.UpdateMovie(r.Context(), id, movie)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, updatedMovie)
}

func (handler *MovieHandler) HandleDeleteMovie(w http.ResponseWriter, r *http.Request) error {
<<<<<<< HEAD
=======
	var req dto.MovieRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	movie := mapper.MovieFromRequest(req)

	createdMovie, err := handler.service.CreateMovie(r.Context(), movie)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusCreated, createdMovie)
}

func (handler *MovieHandler) HandleUpdateMovie(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIDParam(r, "id")
	if err != nil {
		return err
	}

	var req dto.MovieRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	movie := mapper.MovieFromRequest(req)

	updatedMovie, err := handler.service.UpdateMovie(r.Context(), id, movie)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, updatedMovie)
}

func (handler *MovieHandler) HandleDeleteMovie(w http.ResponseWriter, r *http.Request) error {
=======
>>>>>>> da5f31b (feat(movie-service): implement genre management with repository, service, and handler layers; enhance movie handler and routes)
	id, err := parseIDParam(r, "id")
	if err != nil {
		return err
	}

	if err := handler.service.DeleteMovie(r.Context(), id); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
<<<<<<< HEAD
>>>>>>> feature/movieService
=======
>>>>>>> da5f31b (feat(movie-service): implement genre management with repository, service, and handler layers; enhance movie handler and routes)
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

<<<<<<< HEAD
<<<<<<< HEAD
=======
=======
>>>>>>> da5f31b (feat(movie-service): implement genre management with repository, service, and handler layers; enhance movie handler and routes)
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

<<<<<<< HEAD
>>>>>>> feature/movieService
=======
>>>>>>> da5f31b (feat(movie-service): implement genre management with repository, service, and handler layers; enhance movie handler and routes)
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
