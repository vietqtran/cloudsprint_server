package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"go.uber.org/zap"

	"cloud-sprint/config"
	"cloud-sprint/internal/api/router"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/internal/token"

	_ "cloud-sprint/docs/swagger"
)

type Server struct {
	app  *fiber.App
	log  *zap.Logger
	port string
}

func New(store db.Querier, cfg config.Config, log *zap.Logger) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(cfg.JWT.SecretKey, cfg.JWT.RefreshSecretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create token maker: %w", err)
	}

	app := fiber.New(fiber.Config{})

	app.Use(recover.New())
	app.Use(cors.New())

	router.SetupRoutes(app, store, tokenMaker, log, cfg)

	app.Get("/swagger/*", swagger.HandlerDefault)

	return &Server{
		app:  app,
		log:  log,
		port: cfg.Server.Port,
	}, nil
}

func (s *Server) Start(port string) error {
	return s.app.Listen(fmt.Sprintf(":%s", port))
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}