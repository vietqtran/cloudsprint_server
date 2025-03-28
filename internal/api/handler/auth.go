package handler

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"cloud-sprint/config"
	"cloud-sprint/internal/api/request"
	"cloud-sprint/internal/api/response"
	"cloud-sprint/internal/constants"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/internal/token"
	"cloud-sprint/pkg/util"
)

type SetCookieData struct {
	Name      string
	Token     string
	ExpiresAt int
}

type AuthHandler struct {
	store      db.Querier
	tokenMaker token.Maker
	config     config.Config
}

func NewAuthHandler(store db.Querier, tokenMaker token.Maker, config config.Config) *AuthHandler {
	return &AuthHandler{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "Register request"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req request.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.NewErrorResponse(c, constants.StatusBadRequest, "Invalid request body", err).Send(c)
	}

	if err := req.Validate(); err != nil {
		return response.NewErrorResponse(c, constants.StatusBadRequest, err.Error(), nil).Send(c)
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return response.NewErrorResponse(c, constants.StatusBadRequest, "Failed to hash password", err).Send(c)
	}

	user, err := h.store.CreateUser(c.Context(), db.CreateUserParams{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		return response.NewErrorResponse(c, constants.StatusBadRequest, "Failed to create user", err).Send(c)
	}

	_, err = h.store.CreateAccount(c.Context(), db.CreateAccountParams{
		UserID:         user.ID,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return response.InternalServerError(c, "Failed to create account", err)
	}

	return response.Created(c, nil, "User registered successfully")
}

// Login handles user login
// @Summary Login a user
// @Description Login with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "Login request"
// @Success 200 {object} response.LoginResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req request.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	account, err := h.store.GetAccountByEmail(c.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.Unauthorized(c, "Invalid email or password")
		}
		return response.InternalServerError(c, "Failed to get account", err)
	}

	err = util.CheckPassword(req.Password, account.HashedPassword)
	if err != nil {
		return response.Unauthorized(c, "Invalid email or password")
	}

	accessToken, accessPayload, err := h.tokenMaker.CreateToken(
		account.UserID,
		account.Email,
		h.config.JWT.TokenDuration,
	)
	if err != nil {
		return response.InternalServerError(c, "Failed to create access token", err)
	}

	refreshToken, _, err := h.tokenMaker.CreateRefreshToken(
		account.UserID,
		account.Email,
		h.config.JWT.RefreshDuration,
	)
	if err != nil {
		return response.InternalServerError(c, "Failed to create refresh token", err)
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
		return response.InternalServerError(c, "Failed to create session", err)
	}

	SetHttpOnlyCookie(c, SetCookieData{
		Name:      "Refresh",
		Token:     refreshToken,
		ExpiresAt: int(h.config.JWT.RefreshDuration),
	})

	SetHttpOnlyCookie(c, SetCookieData{
		Name:      "Authorization",
		Token:     accessToken,
		ExpiresAt: int(h.config.JWT.TokenDuration),
	})

	user, err := h.store.GetUserByID(c.Context(), account.UserID)
	if err != nil {
		return response.InternalServerError(c, "Failed to get user", err)
	}

	loginResponse := response.NewLoginResponse(user, accessToken, refreshToken, session.ID.String())

	return response.Success(c, loginResponse, "Login successful")
}

// RefreshToken handles token refresh
// @Summary Refresh token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} response.RefreshTokenResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req request.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}
	if err := req.Validate(); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	refreshPayload, err := h.tokenMaker.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		return response.Unauthorized(c, "Invalid refresh token")
	}

	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		return response.BadRequest(c, "Invalid session ID", err)
	}

	session, err := h.store.GetSession(c.Context(), sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.Unauthorized(c, "Session not found")
		}
		return response.InternalServerError(c, "Failed to get session", err)
	}

	if session.IsBlocked {
		return response.Unauthorized(c, "Session is blocked")
	}

	if time.Now().After(session.ExpiresAt) {
		return response.Unauthorized(c, "Session expired")
	}

	if session.RefreshToken != req.RefreshToken {
		return response.Unauthorized(c, "Refresh token doesn't match")
	}

	refreshUserID, err := uuid.Parse(refreshPayload.UserID)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID in refresh token", err)
	}

	account, err := h.store.GetAccountById(c.Context(), session.AccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.Unauthorized(c, "Account not found in session")
		}
		return response.InternalServerError(c, "Failed to get account", err)
	}

	if account.UserID != refreshUserID {
		return response.Unauthorized(c, "User ID doesn't match")
	}

	user, err := h.store.GetUserByID(c.Context(), account.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.NotFound(c, "User not found")
		}
		return response.InternalServerError(c, "Failed to get user", err)
	}

	accessToken, _, err := h.tokenMaker.CreateToken(
		user.ID,
		user.Email,
		h.config.JWT.TokenDuration,
	)
	if err != nil {
		return response.InternalServerError(c, "Failed to create access token", err)
	}

	refreshToken, refreshTokenPayload, err := h.tokenMaker.CreateRefreshToken(
		user.ID,
		user.Email,
		h.config.JWT.RefreshDuration,
	)
	if err != nil {
		return response.InternalServerError(c, "Failed to create refresh token", err)
	}

	_, err = h.store.UpdateSessionRefreshToken(c.Context(), db.UpdateSessionRefreshTokenParams{
		ID:           sessionID,
		RefreshToken: refreshToken,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	})
	if err != nil {
		return response.InternalServerError(c, "Failed to refresh token", err)
	}

	refreshResponse := response.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	SetHttpOnlyCookie(c, SetCookieData{
		Name:      "Refresh",
		Token:     refreshToken,
		ExpiresAt: int(h.config.JWT.RefreshDuration),
	})

	SetHttpOnlyCookie(c, SetCookieData{
		Name:      "Authorization",
		Token:     accessToken,
		ExpiresAt: int(h.config.JWT.TokenDuration),
	})

	return response.Success(c, refreshResponse, "Token refreshed successfully")
}

func SetHttpOnlyCookie(c *fiber.Ctx, data SetCookieData) {
	cookie := new(fiber.Cookie)
	cookie.Name = data.Name
	cookie.Value = data.Token
	cookie.Path = "/"
	cookie.HTTPOnly = true
	cookie.Secure = true
	cookie.MaxAge = 3600 * data.ExpiresAt
	c.Cookie(cookie)
}
