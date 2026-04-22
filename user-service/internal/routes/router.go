package routes

import (
	"net/http"
	"user-service/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterRouter(UserHandler *handlers.UserHandler) http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Get("/{id}", UserHandler.GetUserByID)
		})
	})

	return r
}
