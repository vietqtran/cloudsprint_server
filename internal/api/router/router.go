package router

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"cloud-sprint/config"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/internal/token"
)

func SetupRoutes(app *fiber.App, store db.Querier, tokenMaker token.Maker, logger *zap.Logger, config config.Config, authMiddleware fiber.Handler, refreshMiddleware fiber.Handler) {
	api := app.Group("/api/v1")

	SetupAuthRoutes(api, store, tokenMaker, config, authMiddleware, refreshMiddleware)
	SetupGitHubRoutes(api, store, tokenMaker, config, authMiddleware)
}
