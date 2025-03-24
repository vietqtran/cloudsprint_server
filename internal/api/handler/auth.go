package handler

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"cloud-sprint/config"
	"cloud-sprint/internal/api/request"
	"cloud-sprint/internal/api/response"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/internal/token"
	"cloud-sprint/pkg/util"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	store      db.Querier
	tokenMaker token.Maker
	config     config.Config
}

// NewAuthHandler creates a new auth handler
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
// @Success 201 {object} response.UserResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req request.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
	}

	user, err := h.store.CreateUser(c.Context(), db.CreateUserParams{
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
	}

	userResponse := response.NewUserResponse(user)
	return response.Created(c, userResponse, "User registered successfully")
}

// Login handles user login
// @Summary Login a user
// @Description Login with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "Login request"
// @Success 200 {object} response.LoginResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req request.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	user, err := h.store.GetUser(c.Context(), req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid username or password")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to find user")
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid username or password")
	}

	accessToken, accessPayload, err := h.tokenMaker.CreateToken(
		user.ID,
		user.Username,
		h.config.JWT.TokenDuration,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create access token")
	}

	refreshToken, _, err := h.tokenMaker.CreateToken(
		user.ID,
		user.Username,
		h.config.JWT.TokenDuration*24, // Refresh token lasts 24 times longer
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create refresh token")
	}

	// Create session
	session, err := h.store.CreateSession(c.Context(), db.CreateSessionParams{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    c.Get("User-Agent"),
		ClientIp:     c.IP(),
		ExpiresAt:    accessPayload.ExpiredAt.Add(h.config.JWT.TokenDuration * 24),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create session")
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
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req request.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	// Verify the refresh token
	refreshPayload, err := h.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired refresh token")
	}

	// Check if the session exists
	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid session ID")
	}

	session, err := h.store.GetSession(c.Context(), sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusUnauthorized, "Session not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get session")
	}

	// Check if session is blocked
	if session.IsBlocked {
		return fiber.NewError(fiber.StatusUnauthorized, "Session is blocked")
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		return fiber.NewError(fiber.StatusUnauthorized, "Session has expired")
	}

	// Check if the refresh token matches
	if session.RefreshToken != req.RefreshToken {
		return fiber.NewError(fiber.StatusUnauthorized, "Refresh token doesn't match")
	}

	// Check if the user ID matches
	if session.UserID != refreshPayload.UserID {
		return fiber.NewError(fiber.StatusUnauthorized, "User ID doesn't match")
	}

	// Get the user
	user, err := h.store.GetUserByID(c.Context(), session.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusUnauthorized, "User not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to find user")
	}

	// Create new access token
	accessToken, _, err := h.tokenMaker.CreateToken(
		user.ID,
		user.Username,
		h.config.JWT.TokenDuration,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create access token")
	}

	refreshResponse := response.RefreshTokenResponse{
		AccessToken: accessToken,
	}
	return response.Success(c, refreshResponse, "Token refreshed successfully")
}
