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

	// service
	hallService := service.NewHallService(hallRepo)

	// handlers
	hallHandler := handler.NewHallHandler(hallService)

	// router
	r := routes.RegisterRouter(hallHandler)
	log.Print("Booking Server is running on :8083")
	if err := http.ListenAndServe(":8083", r); err != nil {
		log.Fatal(err)
	}
}
