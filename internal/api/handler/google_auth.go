package handler

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"cloud-sprint/config"
	"cloud-sprint/internal/api/response"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/internal/service"
	"cloud-sprint/internal/token"
	"cloud-sprint/pkg/util"
)

type GoogleAuthHandler struct {
	store        db.Querier
	tokenMaker   token.Maker
	config       config.Config
	emailService *service.EmailService
	googleService *service.GoogleService
}

func NewGoogleAuthHandler(store db.Querier, tokenMaker token.Maker, config config.Config, emailService *service.EmailService, googleService *service.GoogleService) *GoogleAuthHandler {
	return &GoogleAuthHandler{
		store:        store,
		tokenMaker:   tokenMaker,
		config:       config,
		emailService: emailService,
		googleService: googleService,
	}
}

func (h *GoogleAuthHandler) GoogleAuth(c *fiber.Ctx) error {
	oauthConfig := h.getGoogleOAuthConfig()
	url := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return c.Redirect(url)
}

func (h *GoogleAuthHandler) GoogleCallback(c *fiber.Ctx) error {
	fmt.Println("Google callback received")

	code := c.Query("code")
	if code == "" {
		return response.BadRequest(c, "Missing authorization code", nil, nil)
	}

	token, err := h.googleService.Exchange(c.Context(), code)
	if err != nil {
		fmt.Printf("Token exchange error: %v\n", err)
		return response.InternalServerError(c, "Failed to exchange token", err, nil)
	}

	if token == nil {
		return response.InternalServerError(c, "Invalid token received", nil, nil)
	}

	userInfo, err := h.googleService.GetUserInfo(token)
	if err != nil {
		fmt.Printf("User info error: %v\n", err)
		return response.InternalServerError(c, "Failed to get user info", err, nil)
	}

	account, err := h.store.GetAccountByEmail(c.Context(), userInfo.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			return response.InternalServerError(c, "Failed to check account", err, nil)
		}

		user, err := h.store.CreateUser(c.Context(), db.CreateUserParams{
			Email:     userInfo.Email,
			FirstName: userInfo.GivenName,
			LastName:  userInfo.FamilyName,
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

	redirectURL := fmt.Sprintf("%s/auth/callback?access_token=%s&refresh_token=%s&session_id=%s",
		h.config.FrontendBaseURL, accessToken, refreshToken, session.ID.String())
	
	return c.Redirect(redirectURL)
}

func (h *GoogleAuthHandler) getGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     h.config.OAuth.GoogleClientID,
		ClientSecret: h.config.OAuth.GoogleClientSecret,
		RedirectURL:  h.config.OAuth.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
} 