package router

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"cloud-sprint/config"
	"cloud-sprint/internal/api/handler"
	"cloud-sprint/internal/api/middleware"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/internal/token"
)

// SetupRoutes sets up all the routes for the API
func SetupRoutes(app *fiber.App, store db.Querier, tokenMaker token.Maker, logger *zap.Logger, config config.Config) {
	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenMaker)
	loggerMiddleware := middleware.NewLogger(logger)
	errorMiddleware := middleware.NewErrorHandler(logger)

	// Apply global middleware
	app.Use(loggerMiddleware)
	app.Use(errorMiddleware)

	// API routes
	api := app.Group("/api/v1")

	// Auth handler
	authHandler := handler.NewAuthHandler(store, tokenMaker, config)

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// User handler
	userHandler := handler.NewUserHandler(store)

	// User routes (protected by auth middleware)
	users := api.Group("/users", authMiddleware)
	users.Get("/", userHandler.ListUsers)
	users.Get("/me", userHandler.GetCurrentUser)
	users.Get("/:id", userHandler.GetUser)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
}
