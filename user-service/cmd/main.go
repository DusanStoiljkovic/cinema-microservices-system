package main

import (
	"log"
	"net/http"
	"user-service/internal/db"
	"user-service/internal/handlers"
	"user-service/internal/repository"
	"user-service/internal/routes"
	"user-service/internal/services"
)

func main() {

	db, err := db.Connect()
	if err != nil {
		log.Println("Database not connected")
		return
	}

	// repos
	userRepo := repository.NewUserRepository(db)

	// services
	userService := services.NewUserService(userRepo)

	// handlers
	userHandler := handlers.NewUserHandler(userService)

	// router
	r := routes.RegisterRouter(userHandler)

	log.Print("Server is running on :8081")
	http.ListenAndServe(":8081", r)

}
