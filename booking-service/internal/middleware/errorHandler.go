package middleware

import (
	"booking-service/internal/utils"
	"errors"
	"net/http"
)

type AppHandler func(http.ResponseWriter, *http.Request) error

func ErrorHandler(next AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			handleError(w, err)
			return
		}
	}
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, utils.ErrNotFound):
		utils.WriteError(w, http.StatusBadRequest, utils.ErrNotFound.Error())
	case errors.Is(err, utils.ErrConflict):
		utils.WriteError(w, http.StatusConflict, utils.ErrConflict.Error())
	case errors.Is(err, utils.ErrInvalidInput):
		utils.WriteError(w, http.StatusBadRequest, utils.ErrInvalidInput.Error())
	case errors.Is(err, utils.ErrRecordAlreadyExist):
		utils.WriteError(w, http.StatusConflict, utils.ErrRecordAlreadyExist.Error())
	default:
		utils.WriteError(w, http.StatusInternalServerError, "internal server error")
	}

}
