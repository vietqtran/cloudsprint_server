package router

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"cloud-sprint/config"
	"cloud-sprint/internal/api/middleware"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/internal/token"
)

func SetupRoutes(app *fiber.App, store db.Querier, tokenMaker token.Maker, logger *zap.Logger, config config.Config) {
	loggerMiddleware := middleware.NewLogger(logger)
	errorMiddleware := middleware.NewErrorHandler(logger)
	requestIDMiddleware := middleware.RequestID()

	app.Use(middleware.CORSSimple())
	app.Use(requestIDMiddleware)
	app.Use(loggerMiddleware)
	app.Use(errorMiddleware)

	api := app.Group("/api/v1")

	SetupAuthRoutes(api, store, tokenMaker, logger, config)
	SetupUserRoutes(api, store, tokenMaker, logger)
}