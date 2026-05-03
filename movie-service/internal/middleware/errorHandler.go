package middleware

import "net/http"

type AppHandler func(http.ResponseWriter, *http.Request) error

func ErrorHandler(handlerF AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handlerF(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
