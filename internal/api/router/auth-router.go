package router

import (
	"github.com/gofiber/fiber/v2"

	"cloud-sprint/config"
	"cloud-sprint/internal/api/handler"
	_ "cloud-sprint/internal/api/middleware"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/internal/token"
)

func SetupAuthRoutes(api fiber.Router, store db.Querier, tokenMaker token.Maker, config config.Config, authMiddleware fiber.Handler, refreshMiddleware fiber.Handler) {
	authHandler := handler.NewAuthHandler(store, tokenMaker, config)

	auth := api.Group("/auth")
	auth.Post("/sign-up", authHandler.SignUp)
	auth.Post("/sign-in", authHandler.SignIn)
	auth.Post("/refresh", refreshMiddleware, authHandler.RefreshToken)
	auth.Get("/me", authMiddleware, authHandler.Me)
}
