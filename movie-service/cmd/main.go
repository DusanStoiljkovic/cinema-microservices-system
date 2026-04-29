package main

import (
	"log"
	"movie-service/internal/db"
	"movie-service/internal/handler"
	"movie-service/internal/repository"
	"movie-service/internal/routes"
	"movie-service/internal/service"
	"net/http"
)

func main() {
	db, err := db.Connect()
	if err != nil {
		log.Println("Movie Database not connected...")
		return
	}

	// repos
	movieRepo := repository.NewMovieRepository(db)

	// services
	movieService := service.NewMovieService(movieRepo)

	// handlers
	movieHandler := handler.NewMovieHandler(movieService)

	// router
	r := routes.RegisterRouter(movieHandler)

	log.Print("Movie Server is running on :8082")
	http.ListenAndServe(":8082", r)
}
