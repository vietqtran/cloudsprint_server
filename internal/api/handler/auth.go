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

type SetCookieData struct {
	Name      string
	Token     string
	ExpiresAt int
	ENV       string
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

	return response.Created(c, nil, "User signed up successfully! You can sign in now.")
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
		ENV:       h.config.Environment,
	})

	SetHttpOnlyCookie(c, SetCookieData{
		Name:      "Authorization",
		Token:     accessToken,
		ExpiresAt: int(h.config.JWT.TokenDuration),
		ENV:       h.config.Environment,
	})

	user, err := h.store.GetUserByID(c.Context(), account.UserID)
	if err != nil {
		return response.InternalServerError(c, "Failed to get user", err)
	}

	loginResponse := response.NewSignInResponse(user, accessToken, refreshToken, session.ID.String())

	return response.Success(c, loginResponse, "SignIn successful")
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
		return response.BadRequest(c, "invalid request body", err)
	}
	if err := req.Validate(); err != nil {
		return response.BadRequest(c, "invalid request body", err)
	}

	refreshToken, ok := c.Locals("refresh_token").(string)
	if !ok {
		return response.Unauthorized(c, "token is missing")
	}

	refreshPayload, err := h.tokenMaker.VerifyRefreshToken(refreshToken)
	if err != nil {
		return response.Unauthorized(c, "invalid refresh token")
	}

	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		return response.BadRequest(c, "invalid session ID", err)
	}

	session, err := h.store.GetSession(c.Context(), sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.Unauthorized(c, "session not found")
		}
		return response.InternalServerError(c, "Failed to get session", err)
	}

	if session.IsBlocked {
		return response.Unauthorized(c, "session is blocked")
	}

	if time.Now().After(session.ExpiresAt) {
		return response.Unauthorized(c, "session expired")
	}

	if session.RefreshToken != refreshToken {
		return response.Unauthorized(c, "refresh token doesn't match")
	}

	refreshUserID, err := uuid.Parse(refreshPayload.UserID)
	if err != nil {
		return response.BadRequest(c, "invalid user ID in refresh token", err)
	}

	account, err := h.store.GetAccountById(c.Context(), session.AccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.Unauthorized(c, "account not found in session")
		}
		return response.InternalServerError(c, "failed to get account", err)
	}

	if account.UserID != refreshUserID {
		return response.Unauthorized(c, "user id doesn't match")
	}

	user, err := h.store.GetUserByID(c.Context(), account.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.NotFound(c, "user not found")
		}
		return response.InternalServerError(c, "failed to get user", err)
	}

	accessToken, _, err := h.tokenMaker.CreateToken(
		user.ID,
		user.Email,
		h.config.JWT.TokenDuration,
	)
	if err != nil {
		return response.InternalServerError(c, "failed to create access token", err)
	}

	refreshToken, refreshTokenPayload, err := h.tokenMaker.CreateRefreshToken(
		user.ID,
		user.Email,
		h.config.JWT.RefreshDuration,
	)
	if err != nil {
		return response.InternalServerError(c, "failed to create refresh token", err)
	}

	_, err = h.store.UpdateSessionRefreshToken(c.Context(), db.UpdateSessionRefreshTokenParams{
		ID:           sessionID,
		RefreshToken: refreshToken,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	})
	if err != nil {
		return response.InternalServerError(c, "failed to refresh token", err)
	}

	refreshResponse := response.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	SetHttpOnlyCookie(c, SetCookieData{
		Name:      "Refresh",
		Token:     refreshToken,
		ExpiresAt: int(h.config.JWT.RefreshDuration),
		ENV:       h.config.Environment,
	})

	SetHttpOnlyCookie(c, SetCookieData{
		Name:      "Authorization",
		Token:     accessToken,
		ExpiresAt: int(h.config.JWT.TokenDuration),
		ENV:       h.config.Environment,
	})

	return response.Success(c, refreshResponse, "token refreshed successfully")
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID, ok := c.Locals("current_user_id").(string)
	if !ok {
		return response.Unauthorized(c, "user not found")
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
	cookie := new(fiber.Cookie)
	cookie.Name = data.Name
	cookie.Value = data.Token
	cookie.Path = "/"
	cookie.HTTPOnly = true
	cookie.Secure = data.ENV != "development"
	cookie.MaxAge = int(time.Hour) * data.ExpiresAt
	c.Cookie(cookie)
}
