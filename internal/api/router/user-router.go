package router

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"cloud-sprint/internal/api/handler"
	"cloud-sprint/internal/api/middleware"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/internal/token"
)

func SetupUserRoutes(api fiber.Router, store db.Querier, tokenMaker token.Maker, logger *zap.Logger) {
	authMiddleware := middleware.NewAuthMiddleware(tokenMaker)
	userHandler := handler.NewUserHandler(store)

	users := api.Group("/users", authMiddleware)
	users.Get("/", userHandler.ListUsers)
	users.Get("/me", userHandler.GetCurrentUser)
	users.Get("/:id", userHandler.GetUser)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
}