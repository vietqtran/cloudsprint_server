package handler

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"cloud-sprint/config"
	"cloud-sprint/internal/api/request"
	"cloud-sprint/internal/api/response"
	"cloud-sprint/internal/constants"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/internal/service"
	"cloud-sprint/internal/token"
	"cloud-sprint/pkg/util"
)

type SetCookieData struct {
	Name      string
	Token     string
	ExpiresAt int
	ENV       string
}

type AuthHandler struct {
	store        db.Querier
	tokenMaker   token.Maker
	config       config.Config
	emailService *service.EmailService
}

func NewAuthHandler(store db.Querier, tokenMaker token.Maker, config config.Config, emailService *service.EmailService) *AuthHandler {
	return &AuthHandler{
		store:        store,
		tokenMaker:   tokenMaker,
		config:       config,
		emailService: emailService,
	}
}

// SignUp handles user registration
// @Summary SignUp a new user
// @Description SignUp a new user with username, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.SignUpRequest true "SignUp request"
// @Router /auth/sign-up [post]
func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	var req request.SignUpRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return response.BadRequest(c, "Failed to hash password", err)
	}

	_, err = h.store.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			return response.InternalServerError(c, "Failed to get user", err)
		}
	} else {
		return response.BadRequest(c, "User already exists", nil)
	}

	user, err := h.store.CreateUser(c.Context(), db.CreateUserParams{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		return response.BadRequest(c, "Failed to create user", err)
	}

	_, err = h.store.CreateAccount(c.Context(), db.CreateAccountParams{
		UserID:         user.ID,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return response.InternalServerError(c, "Failed to create account", err)
	}

	emailVerificationHandler := NewEmailVerificationHandler(h.store, h.config, h.emailService)
	otp := emailVerificationHandler.generateOTP()
	
	otpID := uuid.New()
	expires := time.Now().Add(15 * time.Minute)

	_, err = h.store.CreateEmailOTP(c.Context(), db.CreateEmailOTPParams{
		ID:        otpID,
		Email:     req.Email,
		Otp:       otp,
		ExpiresAt: expires,
	})
	if err != nil {
		log.Printf("Failed to create OTP: %v", err)
	} else {
		err = h.emailService.SendEmail(service.EmailData{
			To:       req.Email,
			Subject:  "Verify Your Email",
			Template: "email_verification.html",
			Data: map[string]interface{}{
				"Name":      req.Email,
				"OTP":       otp,
				"ExpiresIn": "15 minutes",
			},
		})
		if err != nil {
			log.Printf("Failed to send verification email: %v", err)
		}
	}

	return response.Created(c, nil, "User signed up successfully! Please verify your email.")
}

// SignIn handles user login
// @Summary SignIn a user
// @Description SignIn with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.SignInRequest true "SignIn request"
// @Success 200 {object} response.SignInResponse
// @Router /auth/sign-in [post]
func (h *AuthHandler) SignIn(c *fiber.Ctx) error {
	var req request.SignInRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	account, err := h.store.GetAccountByEmail(c.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.Unauthorized(c, "Invalid email or password", err)
		}
		return response.InternalServerError(c, "Failed to get account", err)
	}

	if !account.EmailVerified {
		return response.Unauthorized(c, "Email not verified", nil)
	}

	err = util.CheckPassword(req.Password, account.HashedPassword)
	if err != nil {
		return response.Unauthorized(c, "Invalid email or password", err)
	}

	accessToken, _, err := h.tokenMaker.CreateToken(
		account.UserID,
		account.Email,
		h.config.JWT.TokenDuration,
	)
	if err != nil {
		return response.InternalServerError(c, "Failed to create access token", err)
	}

	refreshToken, accessPayload, err := h.tokenMaker.CreateRefreshToken(
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
		Name:      "Authorization",
		Token:     accessToken,
		ExpiresAt: int(h.config.JWT.TokenDuration.Seconds()),
		ENV:       h.config.Environment,
	})

	SetHttpOnlyCookie(c, SetCookieData{
		Name:      "Refresh",
		Token:     refreshToken,
		ExpiresAt: int(h.config.JWT.RefreshDuration.Seconds()),
		ENV:       h.config.Environment,
	})

	c.Set("Authorization", "Bearer "+accessToken)

	user, err := h.store.GetUserByID(c.Context(), account.UserID)
	if err != nil {
		return response.InternalServerError(c, "Failed to get user", err)
	}

	loginResponse := response.NewSignInResponse(user, accessToken, refreshToken, session.ID.String())

	return response.Success(c, loginResponse, "SignIn successful")
}

// RefreshToken handles token refresh requests
// @Summary Refresh token
// @Description Refresh access token using refresh token or session ID
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
		return response.BadRequest(c, err.Error(), nil)
	}

	var refreshToken string

	if req.SessionID != "" {
		sessionID, err := uuid.Parse(req.SessionID)
		if err != nil {
			return response.BadRequest(c, "Invalid session ID", err)
		}

		session, err := h.store.GetSession(c.Context(), sessionID)
		if err != nil {
			if err == sql.ErrNoRows {
				return response.NotFound(c, "Session not found", nil)
			}
			return response.InternalServerError(c, "Failed to get session", err)
		}

		if session.IsBlocked {
			return response.Unauthorized(c, "Session is blocked", nil)
		}

		if time.Now().After(session.ExpiresAt) {
			return response.Unauthorized(c, "Session has expired", nil)
		}

		refreshToken = session.RefreshToken
	} else {
		refreshToken = req.RefreshToken

		if refreshToken == "" {
			refreshToken, _ = c.Locals("refresh_token").(string)
		}

		if refreshToken == "" {
			authHeader := c.Get("Refresh")
			if parts := strings.Split(authHeader, " "); len(parts) == 2 && parts[0] == "Bearer" {
				refreshToken = parts[1]
			}
		}
	}

	if refreshToken == "" {
		return response.BadRequest(c, "Refresh token is required", nil)
	}

	refreshPayload, err := h.tokenMaker.VerifyRefreshToken(refreshToken)
	if err != nil {
		return response.Unauthorized(c, "Invalid or expired refresh token", err)
	}

	userUUID, err := uuid.Parse(refreshPayload.UserID)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID format", err)
	}

	accessToken, _, err := h.tokenMaker.CreateToken(
		userUUID,
		refreshPayload.Email,
		h.config.JWT.TokenDuration,
	)
	if err != nil {
		return response.InternalServerError(c, "Failed to create access token", err)
	}

	newRefreshToken, refreshPayload, err := h.tokenMaker.CreateRefreshToken(
		userUUID,
		refreshPayload.Email,
		h.config.JWT.RefreshDuration,
	)
	if err != nil {
		return response.InternalServerError(c, "Failed to create refresh token", err)
	}

	if req.SessionID != "" {
		sessionID, _ := uuid.Parse(req.SessionID)
		_, err = h.store.UpdateSessionRefreshToken(c.Context(), db.UpdateSessionRefreshTokenParams{
			ID:           sessionID,
			RefreshToken: newRefreshToken,
			ExpiresAt:    refreshPayload.ExpiredAt,
		})
		if err != nil {
			return response.InternalServerError(c, "Failed to update session", err)
		}
	}

	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   int(h.config.JWT.TokenDuration.Seconds()),
		Secure:   h.config.Environment == "production",
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "Refresh",
		Value:    newRefreshToken,
		Path:     "/",
		MaxAge:   int(h.config.JWT.RefreshDuration.Seconds()),
		Secure:   h.config.Environment == "production",
		HTTPOnly: true,
	})

	res := response.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}

	return response.NewSuccessResponse(c, constants.StatusOK, res, "Token refreshed successfully").Send(c)
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID, ok := c.Locals("current_user_id").(string)
	if !ok {
		return response.Unauthorized(c, "user not found", nil)
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return response.BadRequest(c, "invalid user id", err)
	}

	user, err := h.store.GetUserByID(c.Context(), userUUID)
	if err != nil {
		return response.InternalServerError(c, "failed to get user", err)
	}

	return response.Success(c, response.NewUserResponse(user), "user found successfully")
}

func SetHttpOnlyCookie(c *fiber.Ctx, data SetCookieData) {
	cookie := &fiber.Cookie{
		Name:     data.Name,
		Value:    data.Token,
		Path:     "/",
		MaxAge:   data.ExpiresAt,
		HTTPOnly: true,
		Secure:   data.ENV == "production",
		SameSite: "Strict",
	}

	c.Cookie(cookie)
}

func (h *AuthHandler) generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(900000) + 100000
	return fmt.Sprintf("%06d", otp)
}
