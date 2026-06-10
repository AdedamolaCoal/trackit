package middleware

import (
	"expense-tracker/pkg/response"
	"log"
	"net/http"
	"runtime/debug"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %v\n%s", err, debug.Stack())
				response.Error(
					w,
					http.StatusInternalServerError,
					response.ErrCodeInternal,
					"An unexpected error occurred",
				)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
