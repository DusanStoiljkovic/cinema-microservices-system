package main

import (
	"api-gateway/internal/proxy"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	p := proxy.NewProxyHandler(5 * time.Second)

	// ROUTES
	p.AddRoute("/api/users", os.Getenv("USER_SERVICE_URL"))
	p.AddRoute("/api/movies", os.Getenv("MOVIE_SERVICE_URL"))
	p.AddRoute("/api/genres", os.Getenv("MOVIE_SERVICE_URL"))
	p.AddRoute("/api/bookings", os.Getenv("BOOKINGS_SERVICE_URL"))
	p.AddRoute("/api/halls", os.Getenv("BOOKINGS_SERVICE_URL"))
	p.AddRoute("/api/projections", os.Getenv("BOOKINGS_SERVICE_URL"))

	log.Println("API Gateway is running on :8080")
	http.ListenAndServe(":8080", p)
}
