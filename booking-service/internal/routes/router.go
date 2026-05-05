package routes

import (
	"booking-service/internal/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HallHander interface {
	HandleGetAllHalls(w http.ResponseWriter, r *http.Request) error
	HandleGetHallByID(w http.ResponseWriter, r *http.Request) error
	HandleCreateHall(w http.ResponseWriter, r *http.Request) error
	HandleUpdateHall(w http.ResponseWriter, r *http.Request) error
	HandleDeleteHall(w http.ResponseWriter, r *http.Request) error
}

type ProjectionHandler interface {
	HandleGetAllProjections(w http.ResponseWriter, r *http.Request) error
	HandleGetProjectionByID(w http.ResponseWriter, r *http.Request) error
	HandleCreateProjection(w http.ResponseWriter, r *http.Request) error
	HandleUpdateProjection(w http.ResponseWriter, r *http.Request) error
	HandleDeleteProjection(w http.ResponseWriter, r *http.Request) error
}

func RegisterRouter(hallHandler HallHander, projectionHandler ProjectionHandler) http.Handler {
	r := chi.NewRouter()

	r.Route("/halls", func(r chi.Router) {
		r.Get("/", middleware.ErrorHandler(hallHandler.HandleGetAllHalls))
		r.Get("/{id}", middleware.ErrorHandler(hallHandler.HandleGetHallByID))
		r.Post("/", middleware.ErrorHandler(hallHandler.HandleCreateHall))
		r.Put("/{id}", middleware.ErrorHandler(hallHandler.HandleUpdateHall))
		r.Delete("/{id}", middleware.ErrorHandler(hallHandler.HandleDeleteHall))
	})

	r.Route("/projections", func(r chi.Router) {
		r.Get("/", middleware.ErrorHandler(projectionHandler.HandleGetAllProjections))
		r.Get("/{id}", middleware.ErrorHandler(projectionHandler.HandleGetProjectionByID))
		r.Post("/", middleware.ErrorHandler(projectionHandler.HandleCreateProjection))
		r.Put("/{id}", middleware.ErrorHandler(projectionHandler.HandleUpdateProjection))
		r.Delete("/{id}", middleware.ErrorHandler(projectionHandler.HandleDeleteProjection))
	})

	return r
}
