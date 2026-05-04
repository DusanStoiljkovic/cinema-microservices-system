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
	genreRepo := repository.NewGenreRepository(db)
	movieRepo := repository.NewMovieRepository(db)

	// services
	genreService := service.NewGenreService(genreRepo)
	movieService := service.NewMovieService(movieRepo, genreRepo)

	// handlers
	genreHandler := handler.NewGenreHandler(genreService)
	movieHandler := handler.NewMovieHandler(movieService)

	// router
	r := routes.RegisterRouter(movieHandler, genreHandler)

	log.Print("Movie Server is running on :8082")
	if err := http.ListenAndServe(":8082", r); err != nil {
		log.Fatal(err)
	}
}
