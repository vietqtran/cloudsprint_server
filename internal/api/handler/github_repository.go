package handler

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"cloud-sprint/config"
	"cloud-sprint/internal/api/response"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/internal/service"
	"cloud-sprint/internal/token"
)

type GitHubRepositoryHandler struct {
	store         db.Querier
	tokenMaker    token.Maker
	config        config.Config
	githubService *service.GitHubService
}

func NewGitHubRepositoryHandler(store db.Querier, tokenMaker token.Maker, config config.Config, githubService *service.GitHubService) *GitHubRepositoryHandler {
	return &GitHubRepositoryHandler{
		store:         store,
		tokenMaker:    tokenMaker,
		config:        config,
		githubService: githubService,
	}
}

// ListRepositories returns all repositories for the authenticated user
// @Summary List GitHub repositories
// @Description Get all GitHub repositories for the authenticated user
// @Tags github
// @Produce json
// @Security BearerAuth
// @Success 200 {array} response.GitHubRepositoryResponse
// @Router /github/repositories [get]
func (h *GitHubRepositoryHandler) ListRepositories(c *fiber.Ctx) error {
	userID, ok := c.Locals("current_user_id").(string)
	if !ok {
		return response.Unauthorized(c, "User not found", nil, nil)
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err, nil)
	}

	account, err := h.store.GetAccountByUserId(c.Context(), userUUID)
	if err != nil {
		return response.InternalServerError(c, "Failed to get account", err, nil)
	}

	oauthAccount, err := h.store.GetOAuthAccountByAccountIDAndProvider(c.Context(), db.GetOAuthAccountByAccountIDAndProviderParams{
		AccountID: account.ID,
		Provider:  "github",
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return response.BadRequest(c, "GitHub account not connected", nil, nil)
		}
		return response.InternalServerError(c, "Failed to get OAuth account", err, nil)
	}

	var token *oauth2.Token
	if oauthAccount.AccessToken.Valid {
		token = &oauth2.Token{
			AccessToken:  oauthAccount.AccessToken.String,
			TokenType:    "Bearer",
			RefreshToken: oauthAccount.RefreshToken.String,
		}

		if oauthAccount.ExpiresAt.Valid {
			token.Expiry = oauthAccount.ExpiresAt.Time
		}

		if token.Expiry.Before(time.Now()) && oauthAccount.RefreshToken.Valid {
			oauthConfig := h.githubService.GetOAuthConfig()
			tokenSource := oauthConfig.TokenSource(c.Context(), token)
			newToken, err := tokenSource.Token()
			if err != nil {
				return response.InternalServerError(c, "Failed to refresh GitHub token", err, nil)
			}
			token = newToken

			_, err = h.store.UpdateOAuthAccount(c.Context(), db.UpdateOAuthAccountParams{
				ID:           oauthAccount.ID,
				AccessToken:  sql.NullString{String: token.AccessToken, Valid: true},
				RefreshToken: sql.NullString{String: token.RefreshToken, Valid: token.RefreshToken != ""},
				ExpiresAt:    sql.NullTime{Time: token.Expiry, Valid: !token.Expiry.IsZero()},
			})
			if err != nil {
				fmt.Printf("Failed to update OAuth account: %v\n", err)
			}
		}
	} else {
		return response.BadRequest(c, "No valid GitHub token found", nil, nil)
	}

	repos, err := h.githubService.GetUserRepositories(token)
	if err != nil {
		return response.InternalServerError(c, "Failed to fetch GitHub repositories", err, nil)
	}

	return response.Success(c, response.NewGitHubRepositoriesResponse(repos), "Repositories retrieved successfully")
}

// GetRepository returns a specific repository by name
// @Summary Get GitHub repository
// @Description Get a specific GitHub repository by name
// @Tags github
// @Produce json
// @Param repo_name path string true "Repository name"
// @Security BearerAuth
// @Success 200 {object} response.GitHubRepositoryResponse
// @Router /github/repositories/{repo_name} [get]
func (h *GitHubRepositoryHandler) GetRepository(c *fiber.Ctx) error {
	repoName := c.Params("repo")
	if repoName == "" {
		return response.BadRequest(c, "Repository name is required", nil, nil)
	}

	userID, ok := c.Locals("current_user_id").(string)
	if !ok {
		return response.Unauthorized(c, "User not found", nil, nil)
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err, nil)
	}

	account, err := h.store.GetAccountByUserId(c.Context(), userUUID)
	if err != nil {
		return response.InternalServerError(c, "Failed to get account", err, nil)
	}

	oauthAccount, err := h.store.GetOAuthAccountByAccountIDAndProvider(c.Context(), db.GetOAuthAccountByAccountIDAndProviderParams{
		AccountID: account.ID,
		Provider:  "github",
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return response.BadRequest(c, "GitHub account not connected", nil, nil)
		}
		return response.InternalServerError(c, "Failed to get OAuth account", err, nil)
	}

	var token *oauth2.Token
	if oauthAccount.AccessToken.Valid {
		token = &oauth2.Token{
			AccessToken:  oauthAccount.AccessToken.String,
			TokenType:    "Bearer",
			RefreshToken: oauthAccount.RefreshToken.String,
		}

		if oauthAccount.ExpiresAt.Valid {
			token.Expiry = oauthAccount.ExpiresAt.Time
		}

		if token.Expiry.Before(time.Now()) && oauthAccount.RefreshToken.Valid {
			oauthConfig := h.githubService.GetOAuthConfig()
			tokenSource := oauthConfig.TokenSource(c.Context(), token)
			newToken, err := tokenSource.Token()
			if err != nil {
				return response.InternalServerError(c, "Failed to refresh GitHub token", err, nil)
			}
			token = newToken

			_, err = h.store.UpdateOAuthAccount(c.Context(), db.UpdateOAuthAccountParams{
				ID:           oauthAccount.ID,
				AccessToken:  sql.NullString{String: token.AccessToken, Valid: true},
				RefreshToken: sql.NullString{String: token.RefreshToken, Valid: token.RefreshToken != ""},
				ExpiresAt:    sql.NullTime{Time: token.Expiry, Valid: !token.Expiry.IsZero()},
			})
			if err != nil {
				fmt.Printf("Failed to update OAuth account: %v\n", err)
			}
		}
	} else {
		return response.BadRequest(c, "No valid GitHub token found", nil, nil)
	}

	repos, err := h.githubService.GetUserRepositories(token)
	if err != nil {
		return response.InternalServerError(c, "Failed to fetch GitHub repositories", err, nil)
	}

	var repo *service.GitHubRepository
	for i, r := range repos {
		if r.Name == repoName {
			repo = &repos[i]
			break
		}
	}

	if repo == nil {
		return response.NotFound(c, "Repository not found", nil, nil)
	}

	return response.Success(c, response.NewGitHubRepositoryResponse(*repo), "Repository retrieved successfully")
}
