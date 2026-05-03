package middleware

import (
	"errors"
	"movie-service/utils"
	"net/http"
)

type AppHandler func(http.ResponseWriter, *http.Request) error

func ErrorHandler(handlerF AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handlerF(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, utils.ErrInvalidInput):
		utils.WriteError(w, http.StatusBadRequest, "invalid input")
	case errors.Is(err, utils.ErrNotFound):
		utils.WriteError(w, http.StatusNotFound, "resource not found")
	case errors.Is(err, utils.ErrConflict):
		utils.WriteError(w, http.StatusConflict, "resource conflict")
	case errors.Is(err, utils.ErrRecordAlreadyExist):
		utils.WriteError(w, http.StatusConflict, "record already exists")
	default:
		utils.WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}
