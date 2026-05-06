package routes

import (
	"net/http"

	"booking-service/internal/middleware"

	"github.com/go-chi/chi/v5"
)

type HallHandler interface {
	HandleGetAllHalls(w http.ResponseWriter, r *http.Request) error
	HandleGetHallByID(w http.ResponseWriter, r *http.Request) error
	HandleCreateHall(w http.ResponseWriter, r *http.Request) error
	HandleUpdateHall(w http.ResponseWriter, r *http.Request) error
	HandleDeleteHall(w http.ResponseWriter, r *http.Request) error
}

type ProjectionHandler interface {
	HandleGetAllProjections(w http.ResponseWriter, r *http.Request) error
	HandleGetProjectionByID(w http.ResponseWriter, r *http.Request) error
	HandleGetProjectionsByMovieID(w http.ResponseWriter, r *http.Request) error
	HandleCreateProjection(w http.ResponseWriter, r *http.Request) error
	HandleUpdateProjection(w http.ResponseWriter, r *http.Request) error
	HandleDeleteProjection(w http.ResponseWriter, r *http.Request) error
}

type TicketHandler interface {
	HandleGetAllTickets(w http.ResponseWriter, r *http.Request) error
	HandleGetTicketByID(w http.ResponseWriter, r *http.Request) error
	HandleGetTicketsByUserID(w http.ResponseWriter, r *http.Request) error
	HandleGetTicketsByProjectionID(w http.ResponseWriter, r *http.Request) error
	HandleCreateTicket(w http.ResponseWriter, r *http.Request) error
	HandleDeleteTicket(w http.ResponseWriter, r *http.Request) error
}

type OrderHandler interface {
	HandleGetAllOrders(w http.ResponseWriter, r *http.Request) error
	HandleGetOrderByID(w http.ResponseWriter, r *http.Request) error
	HandleGetOrdersByUserID(w http.ResponseWriter, r *http.Request) error
	HandleGetMyOrders(w http.ResponseWriter, r *http.Request) error
	HandleCreateOrder(w http.ResponseWriter, r *http.Request) error
	HandlePayOrder(w http.ResponseWriter, r *http.Request) error
	HandleCancelOrder(w http.ResponseWriter, r *http.Request) error
	HandleDeleteOrder(w http.ResponseWriter, r *http.Request) error
}

func RegisterRouter(hallHandler HallHandler, projectionHandler ProjectionHandler, ticketHandler TicketHandler) http.Handler {
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
		r.Get("/movie/{movie_id}", middleware.ErrorHandler(projectionHandler.HandleGetProjectionsByMovieID))
		r.Get("/{id}", middleware.ErrorHandler(projectionHandler.HandleGetProjectionByID))
		r.Post("/", middleware.ErrorHandler(projectionHandler.HandleCreateProjection))
		r.Put("/{id}", middleware.ErrorHandler(projectionHandler.HandleUpdateProjection))
		r.Delete("/{id}", middleware.ErrorHandler(projectionHandler.HandleDeleteProjection))
	})

	r.Route("/tickets", func(r chi.Router) {
		r.Get("/", middleware.ErrorHandler(ticketHandler.HandleGetAllTickets))
		r.Get("/{id}", middleware.ErrorHandler(ticketHandler.HandleGetTicketByID))
		r.Get("/users/{id}", middleware.ErrorHandler(ticketHandler.HandleGetTicketsByUserID))
		r.Get("/projections/{id}", middleware.ErrorHandler(ticketHandler.HandleGetTicketsByProjectionID))
		r.Post("/", middleware.ErrorHandler(ticketHandler.HandleCreateTicket))
		r.Delete("/{id}", middleware.ErrorHandler(ticketHandler.HandleDeleteTicket))
	})

	return r
}
