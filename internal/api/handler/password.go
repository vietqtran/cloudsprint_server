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
	"cloud-sprint/internal/service"
	"cloud-sprint/internal/token"
	"cloud-sprint/pkg/util"
)

type PasswordHandler struct {
	store        db.Querier
	tokenMaker   token.Maker
	config       config.Config
	emailService *service.EmailService
}

func NewPasswordHandler(store db.Querier, tokenMaker token.Maker, config config.Config, emailService *service.EmailService) *PasswordHandler {
	return &PasswordHandler{
		store:        store,
		tokenMaker:   tokenMaker,
		config:       config,
		emailService: emailService,
	}
}

// ForgotPassword handles forgot password requests
// @Summary Request password reset
// @Description Send a password reset email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.ForgotPasswordRequest true "Forgot password request"
// @Success 200 {object} response.BaseResponse
// @Router /auth/forgot-password [post]
func (h *PasswordHandler) ForgotPassword(c *fiber.Ctx) error {
	var req request.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err, nil)
	}

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err.Error(), nil, nil)
	}

	account, err := h.store.GetAccountByEmail(c.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.Success(c, nil, "If your email is registered, you will receive a password reset link")
		}
		return response.InternalServerError(c, "Failed to check email", err, nil)
	}

	resetToken := util.RandomString(64)

	resetID := uuid.New()
	expires := time.Now().Add(time.Hour)

	_, err = h.store.CreatePasswordReset(c.Context(), db.CreatePasswordResetParams{
		ID:        resetID,
		Email:     req.Email,
		Token:     resetToken,
		ExpiresAt: expires,
	})
	if err != nil {
		return response.InternalServerError(c, "Failed to create password reset", err, nil)
	}

	resetURL := fmt.Sprintf("%s/reset?token=%s&email=%s", h.config.FrontendBaseURL, resetToken, req.Email)

	err = h.emailService.SendEmail(service.EmailData{
		To:       req.Email,
		Subject:  "Password Reset",
		Template: "password_reset.html",
		Data: map[string]interface{}{
			"Name":      account.Email,
			"ResetURL":  resetURL,
			"ExpiresIn": "1 hour",
		},
	})
	if err != nil {
		return response.InternalServerError(c, "Failed to send reset email", err, nil)
	}

	return response.Success(c, nil, "If your email is registered, you will receive a password reset link")
}

// VerifyResetToken verifies a password reset token
// @Summary Verify reset token
// @Description Verify if a password reset token is valid
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.VerifyResetTokenRequest true "Verify token request"
// @Success 200 {object} response.BaseResponse
// @Router /auth/verify-reset-token [post]
func (h *PasswordHandler) VerifyResetToken(c *fiber.Ctx) error {
	var req request.VerifyResetTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err, nil)
	}

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err.Error(), nil, nil)
	}

	resetRecord, err := h.store.GetPasswordResetByToken(c.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.BadRequest(c, "Invalid or expired reset token", nil, nil)
		}
		return response.InternalServerError(c, "Failed to verify reset token", err, nil)
	}

	if resetRecord.Token != req.Token {
		return response.BadRequest(c, "Invalid or expired reset token", nil, nil)
	}

	return response.Success(c, nil, "Token is valid")
}

// ResetPassword resets a user's password with a valid token
// @Summary Reset password
// @Description Reset password with a valid token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.ResetPasswordRequest true "Reset password request"
// @Success 200 {object} response.BaseResponse
// @Router /auth/reset-password [post]
func (h *PasswordHandler) ResetPassword(c *fiber.Ctx) error {
	var req request.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err, nil)
	}

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err.Error(), nil, nil)
	}

	resetRecord, err := h.store.GetPasswordResetByToken(c.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.BadRequest(c, "Invalid or expired reset token", nil, nil)
		}
		return response.InternalServerError(c, "Failed to verify reset token", err, nil)
	}

	if resetRecord.Token != req.Token {
		return response.BadRequest(c, "Invalid or expired reset token", nil, nil)
	}

	account, err := h.store.GetAccountByEmail(c.Context(), resetRecord.Email)
	if err != nil {
		return response.InternalServerError(c, "Failed to get account", err, nil)
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return response.InternalServerError(c, "Failed to hash password", err, nil)
	}

	_, err = h.store.UpdateAccountPassword(c.Context(), db.UpdateAccountPasswordParams{
		ID:             account.ID,
		HashedPassword: sql.NullString{String: hashedPassword, Valid: true},
	})
	if err != nil {
		return response.InternalServerError(c, "Failed to update password", err, nil)
	}

	err = h.store.MarkPasswordResetUsed(c.Context(), resetRecord.ID)
	if err != nil {
		return response.InternalServerError(c, "Failed to mark token as used", err, nil)
	}

	return response.Success(c, nil, "Your password has been reset successfully")
}
