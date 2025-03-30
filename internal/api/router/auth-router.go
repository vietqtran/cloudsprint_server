package router

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"cloud-sprint/config"
	"cloud-sprint/internal/api/handler"
	_ "cloud-sprint/internal/api/middleware"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/internal/token"
)

func SetupAuthRoutes(api fiber.Router, store db.Querier, tokenMaker token.Maker, logger *zap.Logger, config config.Config) {
	authHandler := handler.NewAuthHandler(store, tokenMaker, config)

	auth := api.Group("/auth")
	auth.Post("/sign-up", authHandler.SignUp)
	auth.Post("/sign-in", authHandler.SignIn)
	auth.Post("/refresh", authHandler.RefreshToken)
}