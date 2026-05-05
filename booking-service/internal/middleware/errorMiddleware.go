package middleware

import (
	"booking-service/internal/utils"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) error

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ErrorHandler(handler AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			HandleHTTPError(w, r, err)
			return
		}
	}
}

func HandleHTTPError(w http.ResponseWriter, r *http.Request, err error) {
	var safeErr *utils.SafeError

	if !errors.As(err, &safeErr) {
		safeErr = utils.NewInternal("Internal server error", err)
	}

	slog.Error(
		"request failed",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.Any("error", safeErr),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(safeErr.Status)

	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Code:    safeErr.Code,
		Message: safeErr.UserMsg,
	})
}
