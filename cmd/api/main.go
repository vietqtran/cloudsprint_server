package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"cloud-sprint/config"
	"cloud-sprint/internal/api/server"
	"cloud-sprint/internal/db"
	"cloud-sprint/internal/logger"
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
// @in header
// @name Authorization
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load configuration: %w", err))
	}

	log, err := logger.NewLogger(cfg.Environment)
	if err != nil {
		panic(fmt.Errorf("failed to initialize logger: %w", err))
	}
	defer log.Sync()

	conn, queries, err := db.Connect(cfg.Database, log)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer conn.Close()

	app, err := server.New(queries, cfg, log)
	if err != nil {
		log.Fatal("failed to create server", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := app.Start(cfg.Server.Port); err != nil {
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