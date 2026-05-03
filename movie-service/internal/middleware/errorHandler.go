package middleware

import "net/http"

type AppHandler func(http.ResponseWriter, *http.Request) error

<<<<<<< HEAD
<<<<<<< HEAD
func ErrorHandler(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
=======
func ErrorHandler(handlerF AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handlerF(w, r)
>>>>>>> feature/movieService
=======
func ErrorHandler(handlerF AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handlerF(w, r)
>>>>>>> da5f31b (feat(movie-service): implement genre management with repository, service, and handler layers; enhance movie handler and routes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
