package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"movie-service/internal/utils"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) error

func ErrorHandler(next AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			handleError(w, r, err)
			return
		}
	}
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
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

	utils.WriteJSON(w, safeErr.Status, map[string]string{
		"error": safeErr.UserMsg,
		"code":  safeErr.Code,
	})
}
