package categories

import (
	"expense-tracker/internal/shared/middleware"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, auth *middleware.AuthMiddleware) {
	mux.Handle("GET /api/v1/categories",
		auth.Authenticate(http.HandlerFunc(handler.GetAll)))

	mux.Handle("POST /api/v1/categories",
		auth.Authenticate(http.HandlerFunc(handler.Create)))

	mux.Handle("GET /api/v1/categories/{id}",
		auth.Authenticate(http.HandlerFunc(handler.GetByID)))

	mux.Handle("PUT /api/v1/categories/{id}",
		auth.Authenticate(http.HandlerFunc(handler.Update)))

	mux.Handle("DELETE /api/v1/categories/{id}",
		auth.Authenticate(http.HandlerFunc(handler.Delete)))
}
