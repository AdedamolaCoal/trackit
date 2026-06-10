package main

import (
	"expense-tracker/internal/config"
	"expense-tracker/internal/database"
	"expense-tracker/internal/modules/auth"
	"expense-tracker/internal/modules/categories"
	"expense-tracker/internal/shared/middleware"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()
	cfg.Validate()

	db := database.Connect(cfg)

	database.Migrate(db.DB, &auth.User{}, &auth.RefreshToken{}, &categories.Category{})

	mux := http.NewServeMux()

	// Module Wiring
	// Auth
	authRepo := auth.NewRepository(db.DB)
	authSrv := auth.NewService(cfg, authRepo)
	authHandler := auth.NewHandler(authSrv)
	authMiddleware := middleware.NewAuthMiddleware(cfg)

	// Categories
	categoryRepo := categories.NewRepository(db.DB)
	categorySrv := categories.NewService(categoryRepo)
	categoryHandler := categories.NewHandler(categorySrv)

	//	Routes
	auth.RegisterRoutes(mux, authHandler, authMiddleware)
	categories.RegisterRoutes(mux, categoryHandler, authMiddleware)

	//	Middleware Stack - outermost runs first
	stack := middleware.Recover(
		middleware.Logger(
			middleware.CORS(mux),
		),
	)

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Printf("Server starting on %s: (env: %s)", addr, cfg.App.Env)

	if err := http.ListenAndServe(addr, stack); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
