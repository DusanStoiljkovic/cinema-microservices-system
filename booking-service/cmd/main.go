package main

import (
	"booking-service/internal/db"
	"booking-service/internal/handler"
	"booking-service/internal/repository"
	"booking-service/internal/routes"
	"booking-service/internal/service"
	"log"
	"net/http"
)

func main() {
	db, err := db.Connect()
	if err != nil {
		log.Println("Booking Database not connected...")
		return
	}

	// repos
	hallRepo := repository.NewHallRepository(db)
	projectionRepo := repository.NewProjectionRepository(db)
	ticketRepo := repository.NewTicketRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// service
	hallService := service.NewHallService(hallRepo)
	projectionService := service.NewProjectionService(projectionRepo)
	ticketService := service.NewTicketService(ticketRepo)
	orderService := service.NewOrderService(orderRepo)

	// handlers
	hallHandler := handler.NewHallHandler(hallService)
	projectionHandler := handler.NewProjectionHandler(projectionService)
	ticketHandler := handler.NewTicketHandler(ticketService)
	orderHandler := handler.NewOrderHandler(orderService)

	// router
	r := routes.RegisterRouter(hallHandler, projectionHandler, ticketHandler, orderHandler)
	log.Print("Booking Server is running on :8083")
	if err := http.ListenAndServe(":8083", r); err != nil {
		log.Fatal(err)
	}
}
