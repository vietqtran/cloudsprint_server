package router

import (
	"github.com/gofiber/fiber/v2"

	"cloud-sprint/config"
	"cloud-sprint/internal/api/handler"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/internal/service"
	"cloud-sprint/internal/token"
)

func SetupGitHubRoutes(api fiber.Router, store db.Querier, tokenMaker token.Maker, config config.Config, authMiddleware fiber.Handler) {
	githubService := service.NewGitHubService(config)
	
	githubHandler := handler.NewGitHubRepositoryHandler(store, tokenMaker, config, githubService)
	
	github := api.Group("/github")
	github.Get("/repositories", authMiddleware, githubHandler.ListRepositories)
	github.Get("/repository/:owner/:repo", authMiddleware, githubHandler.GetRepository)
}
