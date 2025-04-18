package router

import (
	"github.com/gofiber/fiber/v2"

	"cloud-sprint/config"
	"cloud-sprint/internal/api/handler"
	"cloud-sprint/internal/api/middleware"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/internal/service"
	"cloud-sprint/internal/token"
)

func SetupGitHubRoutes(api fiber.Router, store db.Querier, tokenMaker token.Maker, config config.Config) {
	githubService := service.NewGitHubService(config)

	githubRepoHandler := handler.NewGitHubRepositoryHandler(store, tokenMaker, config, githubService)

	authMiddleware := middleware.NewAuthMiddleware(tokenMaker, "access", store, config)

	github := api.Group("/github")
	github.Use(authMiddleware)

	repositories := github.Group("/repositories")
	repositories.Get("/", githubRepoHandler.ListRepositories)
	repositories.Get("/:repo_name", githubRepoHandler.GetRepository)
}
