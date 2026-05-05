package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"movie-service/internal/dto"
	"movie-service/internal/models"
	"movie-service/internal/utils"
)

type GenreService interface {
	GetGenres(ctx context.Context) ([]*models.Genre, error)
	GetGenreByFilter(ctx context.Context, req *dto.GenreFilter) (*models.Genre, error)
	CreateGenre(ctx context.Context, genre *models.Genre) (*models.Genre, error)
	UpdateGenre(ctx context.Context, genre *models.Genre, id uint) (*models.Genre, error)
	DeleteGenre(ctx context.Context, id uint) error
}

type GenreHandler struct {
	service GenreService
}

func NewGenreHandler(service GenreService) *GenreHandler {
	return &GenreHandler{service: service}
}

func (handler *GenreHandler) HandleGetGenres(w http.ResponseWriter, r *http.Request) error {
	genres, err := handler.service.GetGenres(r.Context())
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, genres)
}

func (handler *GenreHandler) HandleGetGenreByFilter(w http.ResponseWriter, r *http.Request) error {
	req := &dto.GenreFilter{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.NewInvalidInput("Invalid request body", err)
	}

	genre, err := handler.service.GetGenreByFilter(r.Context(), req)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, genre)
}

func (handler *GenreHandler) HandleCreateGenre(w http.ResponseWriter, r *http.Request) error {
	req := &dto.GenreRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.NewInvalidInput("Invalid request body", err)
	}

	genre := &models.Genre{
		Name: req.Name,
	}

	createdGenre, err := handler.service.CreateGenre(r.Context(), genre)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusCreated, createdGenre)
}

func (handler *GenreHandler) HandleUpdateGenre(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIDParam(r, "id")
	if err != nil {
		return utils.NewInvalidInput("Invalid genre id", err)
	}

	req := &dto.GenreRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.NewInvalidInput("Invalid request body", err)
	}

	genre := &models.Genre{
		Name: req.Name,
	}

	updatedGenre, err := handler.service.UpdateGenre(r.Context(), genre, id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, updatedGenre)
}

func (handler *GenreHandler) HandleDeleteGenre(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIDParam(r, "id")
	if err != nil {
		return utils.NewInvalidInput("Invalid genre id", err)
	}

	if err := handler.service.DeleteGenre(r.Context(), id); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
