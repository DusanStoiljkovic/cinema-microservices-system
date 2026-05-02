package middleware

import "net/http"

type AppHandler func(http.ResponseWriter, *http.Request) error

<<<<<<< HEAD
func ErrorHandler(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
=======
func ErrorHandler(handlerF AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handlerF(w, r)
>>>>>>> feature/movieService
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
