package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"user-service/internal/secure"
	"user-service/internal/utils"
)

func ErrorHandler(next AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := next(w, r)
		if err == nil {
			return
		}

		var safeErr *secure.SafeError

		if errors.As(err, &safeErr) {
			fmt.Println(safeErr.LogString()) // log

			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
				"error": safeErr.UserMsg,
				"code":  safeErr.Code,
			})
			return
		}

		fmt.Println("UNKNOWN ERROR: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
