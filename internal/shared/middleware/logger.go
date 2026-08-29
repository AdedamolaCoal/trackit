package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes = n
	return n, err
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{w, http.StatusOK, 0}
		next.ServeHTTP(wrapped, r)

		log.Printf(
			"%s %s %d %dB %s",
			r.Method,
			r.URL.Path,
			wrapped.status,
			wrapped.bytes,
			time.Since(start),
		)
	})
}
