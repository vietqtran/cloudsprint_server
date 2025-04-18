package router

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"cloud-sprint/config"
	"cloud-sprint/internal/api/middleware"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/internal/token"
)

func SetupRoutes(app *fiber.App, store db.Querier, tokenMaker token.Maker, logger *zap.Logger, config config.Config, authMiddleware fiber.Handler, refreshMiddleware fiber.Handler) {
	api := app.Group("/api/v1")

	authMiddleware = middleware.NewAuthMiddleware(tokenMaker, "access", store, config)
	refreshMiddleware = middleware.NewAuthMiddleware(tokenMaker, "refresh", store, config)

	SetupAuthRoutes(api, store, tokenMaker, config)
	SetupGitHubRoutes(api, store, tokenMaker, config)
}
