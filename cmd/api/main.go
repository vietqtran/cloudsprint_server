package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"cloud-sprint/config"
	"cloud-sprint/internal/api/router"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/internal/logger"
	"cloud-sprint/internal/token"

	_ "cloud-sprint/docs/swagger"
)

// @title Go Postgres API
// @version 1.0
// @description A RESTful API built with Go, Fiber, and PostgreSQL
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email your.email@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.LoadConfig("config")
	if err != nil {
		panic(fmt.Errorf("failed to load configuration: %w", err))
	}

	log, err := logger.NewLogger(cfg.Environment)
	if err != nil {
		panic(fmt.Errorf("failed to initialize logger: %w", err))
	}
	defer log.Sync()

	conn, err := sql.Open(cfg.Database.Driver, cfg.Database.Source)
	if err != nil {
		log.Fatal("cannot connect to database", zap.Error(err))
	}
	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		log.Fatal("cannot ping database", zap.Error(err))
	}
	log.Info("database connected successfully")

	tokenMaker, err := token.NewJWTMaker(cfg.JWT.SecretKey, cfg.JWT.RefreshSecretKey)
	if err != nil {
		log.Fatal("cannot create token maker", zap.Error(err))
	}

	queries := db.New(conn)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			}

			log.Error("error handler caught error",
				zap.Int("status", code),
				zap.String("message", message),
				zap.Error(err),
			)

			return c.Status(code).JSON(fiber.Map{
				"error": message,
			})
		},
	})

	app.Use(recover.New())
	app.Use(cors.New())

	router.SetupRoutes(app, queries, tokenMaker, log, cfg)

	app.Get("/swagger/*", swagger.HandlerDefault)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
			log.Fatal("error starting server", zap.Error(err))
		}
	}()

	log.Info(fmt.Sprintf("server started on http://localhost:%s", cfg.Server.Port))
	log.Info(fmt.Sprintf("swagger UI available at http://localhost:%s/swagger/", cfg.Server.Port))

	<-ctx.Done()
	log.Info("shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Fatal("error shutting down server", zap.Error(err))
	}

	log.Info("server gracefully stopped")
}
