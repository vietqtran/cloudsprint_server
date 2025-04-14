package handler

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"cloud-sprint/config"
	"cloud-sprint/internal/api/response"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/internal/service"
	"cloud-sprint/internal/token"
	"cloud-sprint/pkg/util"
)

type GitHubAuthHandler struct {
	store         db.Querier
	tokenMaker    token.Maker
	config        config.Config
	emailService  *service.EmailService
	githubService *service.GitHubService
}

func NewGitHubAuthHandler(store db.Querier, tokenMaker token.Maker, config config.Config, emailService *service.EmailService, githubService *service.GitHubService) *GitHubAuthHandler {
	return &GitHubAuthHandler{
		store:         store,
		tokenMaker:    tokenMaker,
		config:        config,
		emailService:  emailService,
		githubService: githubService,
	}
}

// GitHubAuth initiates the GitHub OAuth flow
// @Summary Initiate GitHub OAuth
// @Description Redirect to GitHub for authentication
// @Tags auth
// @Produce json
// @Router /auth/github/auth [get]
func (h *GitHubAuthHandler) GitHubAuth(c *fiber.Ctx) error {
	state := util.RandomString(16)

	sessionID := uuid.New().String()
	stateKey := fmt.Sprintf("github_oauth_state:%s", sessionID)

	c.Cookie(&fiber.Cookie{
		Name:     "github_session",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   int(15 * time.Minute.Seconds()),
		Secure:   h.config.Environment == "production",
		HTTPOnly: true,
	})

	fmt.Printf("Stored state %s for session %s\n", state, stateKey)

	oauthConfig := h.githubService.GetOAuthConfig()
	url := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
	return c.Redirect(url)
}

// GitHubCallback processes the GitHub OAuth callback
// @Summary GitHub OAuth callback
// @Description Process the callback from GitHub OAuth
// @Tags auth
// @Produce json
// @Param code query string true "Authorization code"
// @Param state query string true "State for CSRF protection"
// @Router /auth/github/callback [get]
func (h *GitHubAuthHandler) GitHubCallback(c *fiber.Ctx) error {
	fmt.Println("GitHub callback received")

	code := c.Query("code")
	if code == "" {
		return response.BadRequest(c, "Missing authorization code", nil, nil)
	}

	// In production, validate the state parameter
	// receivedState := c.Query("state")
	// sessionID := c.Cookies("github_session")
	// stateKey := fmt.Sprintf("github_oauth_state:%s", sessionID)
	// expectedState := getStateFromStorage(stateKey)
	// if receivedState != expectedState {
	//     return response.BadRequest(c, "Invalid state parameter", nil, nil)
	// }

	// Exchange code for token
	token, err := h.githubService.Exchange(c.Context(), code)
	if err != nil {
		fmt.Printf("Token exchange error: %v\n", err)
		return response.InternalServerError(c, "Failed to exchange token", err, nil)
	}

	if token == nil {
		return response.InternalServerError(c, "Invalid token received", nil, nil)
	}

	userInfo, err := h.githubService.GetUserInfo(token)
	if err != nil {
		fmt.Printf("User info error: %v\n", err)
		return response.InternalServerError(c, "Failed to get user info", err, nil)
	}

	if userInfo.Email == "" {
		return response.BadRequest(c, "GitHub account does not have a verified email", nil, nil)
	}

	account, err := h.store.GetAccountByEmail(c.Context(), userInfo.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			return response.InternalServerError(c, "Failed to check account", err, nil)
		}

		firstName, lastName := parseFullName(userInfo.Name)
		if firstName == "" {
			firstName = userInfo.Login 
		}
		if lastName == "" {
			lastName = "-"
		}

		user, err := h.store.CreateUser(c.Context(), db.CreateUserParams{
			Email:     userInfo.Email,
			FirstName: firstName,
			LastName:  lastName,
		})
		if err != nil {
			return response.InternalServerError(c, "Failed to create user", err, nil)
		}

		randomPassword := util.RandomPassword()
		hashedPassword, err := util.HashPassword(randomPassword)
		if err != nil {
			return response.InternalServerError(c, "Failed to hash password", err, nil)
		}

		account, err = h.store.CreateAccount(c.Context(), db.CreateAccountParams{
			UserID:         user.ID,
			Email:          userInfo.Email,
			HashedPassword: sql.NullString{String: hashedPassword, Valid: true},
		})
		if err != nil {
			return response.InternalServerError(c, "Failed to create account", err, nil)
		}

		_, err = h.store.UpdateAccountEmailVerificationStatus(c.Context(), db.UpdateAccountEmailVerificationStatusParams{
			ID:            account.ID,
			EmailVerified: true,
		})
		if err != nil {
			fmt.Printf("Failed to verify email: %v\n", err)
		}

		_, err = h.store.CreateOAuthAccount(c.Context(), db.CreateOAuthAccountParams{
			AccountID:      account.ID,
			Provider:       "github",
			ProviderUserID: fmt.Sprintf("%d", userInfo.ID),
			AccessToken:    sql.NullString{String: token.AccessToken, Valid: true},
			RefreshToken:   sql.NullString{String: token.RefreshToken, Valid: token.RefreshToken != ""},
			ExpiresAt:      sql.NullTime{Time: token.Expiry, Valid: !token.Expiry.IsZero()},
		})
		if err != nil {
			fmt.Printf("Failed to create OAuth account: %v\n", err)
		}
	} else {
		existingOAuth, err := h.store.GetOAuthAccountByProviderAndProviderUserID(c.Context(), db.GetOAuthAccountByProviderAndProviderUserIDParams{
			Provider:       "github",
			ProviderUserID: fmt.Sprintf("%d", userInfo.ID),
		})

		if err != nil {
			if err == sql.ErrNoRows {
				_, err = h.store.CreateOAuthAccount(c.Context(), db.CreateOAuthAccountParams{
					AccountID:      account.ID,
					Provider:       "github",
					ProviderUserID: fmt.Sprintf("%d", userInfo.ID),
					AccessToken:    sql.NullString{String: token.AccessToken, Valid: true},
					RefreshToken:   sql.NullString{String: token.RefreshToken, Valid: token.RefreshToken != ""},
					ExpiresAt:      sql.NullTime{Time: token.Expiry, Valid: !token.Expiry.IsZero()},
				})
				if err != nil {
					fmt.Printf("Failed to create OAuth account: %v\n", err)
				}
			} else {
				return response.InternalServerError(c, "Failed to check OAuth account", err, nil)
			}
		} else {
			_, err = h.store.UpdateOAuthAccount(c.Context(), db.UpdateOAuthAccountParams{
				ID:           existingOAuth.ID,
				AccessToken:  sql.NullString{String: token.AccessToken, Valid: true},
				RefreshToken: sql.NullString{String: token.RefreshToken, Valid: token.RefreshToken != ""},
				ExpiresAt:    sql.NullTime{Time: token.Expiry, Valid: !token.Expiry.IsZero()},
			})
			if err != nil {
				fmt.Printf("Failed to update OAuth account: %v\n", err)
			}
		}
	}

	accessToken, _, err := h.tokenMaker.CreateToken(
		account.UserID,
		account.Email,
		h.config.JWT.TokenDuration,
	)
	if err != nil {
		return response.InternalServerError(c, "Failed to create access token", err, nil)
	}

	refreshToken, accessPayload, err := h.tokenMaker.CreateRefreshToken(
		account.ID,
		account.Email,
		h.config.JWT.RefreshDuration,
	)
	if err != nil {
		return response.InternalServerError(c, "Failed to create refresh token", err, nil)
	}

	session, err := h.store.CreateSession(c.Context(), db.CreateSessionParams{
		ID:           uuid.New(),
		AccountID:    account.ID,
		RefreshToken: refreshToken,
		UserAgent:    c.Get("User-Agent"),
		ClientIp:     c.IP(),
		ExpiresAt:    accessPayload.ExpiredAt.Add(h.config.JWT.TokenDuration),
	})
	if err != nil {
		return response.InternalServerError(c, "Failed to create session", err, nil)
	}

	util.SetHttpOnlyCookie(c, util.SetCookieData{
		Name:      "Authorization",
		Token:     accessToken,
		ExpiresAt: int(h.config.JWT.TokenDuration.Seconds()),
		ENV:       h.config.Environment,
	})

	util.SetHttpOnlyCookie(c, util.SetCookieData{
		Name:      "Refresh",
		Token:     refreshToken,
		ExpiresAt: int(h.config.JWT.RefreshDuration.Seconds()),
		ENV:       h.config.Environment,
	})

	redirectURL := fmt.Sprintf("%s/auth/callback?access_token=%s&refresh_token=%s&session_id=%s&provider=github",
		h.config.FrontendBaseURL, accessToken, refreshToken, session.ID.String())

	return c.Redirect(redirectURL)
}

func parseFullName(fullName string) (string, string) {
	if fullName == "" {
		return "", ""
	}

	parts := strings.Split(fullName, " ")
	if len(parts) == 1 {
		return parts[0], ""
	}

	return parts[0], parts[len(parts)-1]
}
