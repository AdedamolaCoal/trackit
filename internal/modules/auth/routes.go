package auth

import (
	"expense-tracker/internal/shared/middleware"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, auth *middleware.AuthMiddleware) {
	//	Public Routes
	mux.HandleFunc("POST /api/v1/auth/register", handler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", handler.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", handler.Refresh)
	mux.HandleFunc("POST /api/v1/auth/logout", handler.Logout)

	//	 Protected Routes
	mux.Handle("GET /api/v1/users/me", auth.Authenticate(http.HandlerFunc(handler.GetProfile)))
	mux.Handle("PUT /api/v1/users/me", auth.Authenticate(http.HandlerFunc(handler.UpdateProfile)))
	mux.Handle("PUT /api/v1/users/me/password", auth.Authenticate(http.HandlerFunc(handler.ChangePassword)))
}
