package handler

import (
	"database/sql"
	"fmt"
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
		return response.BadRequest(c, "Failed to hash password", nil)
	}

	user, err := h.store.CreateUser(c.Context(), db.CreateUserParams{
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
	})
	if err != nil {
		return response.BadRequest(c, "Failed to create user", nil)
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
			return response.Unauthorized(c, "Invalid username or password")
		}
		return response.InternalServerError(c, "Failed to get user")
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return response.Unauthorized(c, "Invalid username or password")
	}

	accessToken, accessPayload, err := h.tokenMaker.CreateToken(
		user.ID,
		user.Username,
		h.config.JWT.TokenDuration,
	)
	if err != nil {
		return response.InternalServerError(c, "Failed to create access token")
	}

	refreshToken, _, err := h.tokenMaker.CreateRefreshToken(
		user.ID,
		user.Username,
		h.config.JWT.TokenDuration*24,
	)
	if err != nil {
		return response.InternalServerError(c, "Failed to create refresh token")
	}

	session, err := h.store.CreateSession(c.Context(), db.CreateSessionParams{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    c.Get("User-Agent"),
		ClientIp:     c.IP(),
		ExpiresAt:    accessPayload.ExpiredAt.Add(h.config.JWT.TokenDuration * 24),
	})
	if err != nil {
		return response.InternalServerError(c, "Failed to create session")
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

	refreshPayload, err := h.tokenMaker.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		return response.Unauthorized(c, "Invalid refresh token")
	}

	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		return response.BadRequest(c, "Invalid session ID", nil)
	}

	session, err := h.store.GetSession(c.Context(), sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.Unauthorized(c, "Session not found")
		}
		return response.InternalServerError(c, "Failed to get session")
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
		return response.BadRequest(c, "Invalid user ID in refresh token", nil)
	}
	if session.UserID != refreshUserID {
		return response.Unauthorized(c, "User ID doesn't match")
	}

	user, err := h.store.GetUserByID(c.Context(), session.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.NotFound(c, "User not found")
		}
		return response.InternalServerError(c, "Failed to get user")
	}

	accessToken, _, err := h.tokenMaker.CreateToken(
		user.ID,
		user.Username,
		h.config.JWT.TokenDuration,
	)
	if err != nil {
		return response.InternalServerError(c, "Failed to create access token")
	}

	refreshToken, refreshTokenPayload, err := h.tokenMaker.CreateRefreshToken(
		user.ID,
		user.Username,
		h.config.JWT.RefreshDuration,
	)
	if err != nil {
		return response.InternalServerError(c, "Failed to create refresh token")
	}

	updatedSession, err := h.store.UpdateSessionRefreshToken(c.Context(), db.UpdateSessionRefreshTokenParams{
		ID:           sessionID,
		RefreshToken: refreshToken,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	})
	if err != nil {
		fmt.Println("Error updating session:", err)
		return response.InternalServerError(c, "Failed to refresh token")
	}
	fmt.Println("Session updated successfully:", updatedSession.ID)

	refreshResponse := response.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return response.Success(c, refreshResponse, "Token refreshed successfully")
}