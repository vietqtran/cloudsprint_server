package handler

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"cloud-sprint/config"
	"cloud-sprint/internal/api/request"
	"cloud-sprint/internal/api/response"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/internal/service"
)

type EmailVerificationHandler struct {
	store        db.Querier
	config       config.Config
	emailService *service.EmailService
}

func NewEmailVerificationHandler(store db.Querier, config config.Config, emailService *service.EmailService) *EmailVerificationHandler {
	return &EmailVerificationHandler{
		store:        store,
		config:       config,
		emailService: emailService,
	}
}

func (h *EmailVerificationHandler) generateOTP() string {
	otp := rand.Intn(900000) + 100000
	return fmt.Sprintf("%06d", otp)
}

// SendOTP handles sending email verification OTP
// @Summary Send email verification OTP
// @Description Send a one-time password to verify email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.SendEmailOTPRequest true "Send OTP request"
// @Success 200 {object} response.BaseResponse
// @Router /auth/verify-email/send-otp [post]
func (h *EmailVerificationHandler) SendOTP(c *fiber.Ctx) error {
	var req request.SendEmailOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	account, err := h.store.GetAccountByEmail(c.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.NotFound(c, "Account not found", nil)
		}
		return response.InternalServerError(c, "Failed to check email", err)
	}

	otp := h.generateOTP()
	
	otpID := uuid.New()
	expires := time.Now().Add(15 * time.Minute)

	_, err = h.store.CreateEmailOTP(c.Context(), db.CreateEmailOTPParams{
		ID:        otpID,
		Email:     req.Email,
		Otp:       otp,
		ExpiresAt: expires,
	})
	if err != nil {
		return response.InternalServerError(c, "Failed to create OTP", err)
	}

	err = h.emailService.SendEmail(service.EmailData{
		To:       req.Email,
		Subject:  "Email Verification",
		Template: "email_verification.html",
		Data: map[string]interface{}{
			"Name":      account.Email,
			"OTP":       otp,
			"ExpiresIn": "15 minutes",
		},
	})
	if err != nil {
		return response.InternalServerError(c, "Failed to send verification email", err)
	}

	return response.Success(c, nil, "Verification code sent to your email")
}

// VerifyOTP verifies an email OTP
// @Summary Verify email with OTP
// @Description Verify email address using OTP
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.VerifyEmailOTPRequest true "Verify OTP request"
// @Success 200 {object} response.EmailVerificationResponse
// @Router /auth/verify-email/verify [post]
func (h *EmailVerificationHandler) VerifyOTP(c *fiber.Ctx) error {
	var req request.VerifyEmailOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	if err := req.Validate(); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	otpRecord, err := h.store.GetEmailOTPByCode(c.Context(), db.GetEmailOTPByCodeParams{
		Otp:   req.OTP,
		Email: req.Email,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return response.BadRequest(c, "Invalid or expired verification code", nil)
		}
		return response.InternalServerError(c, "Failed to verify code", err)
	}

	account, err := h.store.GetAccountByEmail(c.Context(), req.Email)
	if err != nil {
		return response.InternalServerError(c, "Failed to get account", err)
	}

	err = h.store.MarkEmailOTPUsed(c.Context(), otpRecord.ID)
	if err != nil {
		return response.InternalServerError(c, "Failed to mark code as used", err)
	}

	updatedAccount, err := h.store.UpdateAccountEmailVerified(c.Context(), account.ID)
	if err != nil {
		return response.InternalServerError(c, "Failed to update verification status", err)
	}

	return response.Success(c, response.EmailVerificationResponse{
		EmailVerified: updatedAccount.EmailVerified,
	}, "Email verified successfully")
}

// CheckVerificationStatus checks if a user's email is verified
// @Summary Check email verification status
// @Description Check if the authenticated user's email is verified
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.EmailVerificationResponse
// @Router /auth/verify-email/status [get]
func (h *EmailVerificationHandler) CheckVerificationStatus(c *fiber.Ctx) error {
	accountID, ok := c.Locals("current_account_id").(string)
	if !ok {
		return response.Unauthorized(c, "account not found", nil)
	}

	accountUUID, err := uuid.Parse(accountID)
	if err != nil {
		return response.BadRequest(c, "invalid account id", err)
	}

	emailVerified, err := h.store.GetEmailVerificationStatus(c.Context(), accountUUID)
	if err != nil {
		return response.InternalServerError(c, "Failed to get verification status", err)
	}

	return response.Success(c, response.EmailVerificationResponse{
		EmailVerified: emailVerified,
	}, "Verification status retrieved")
} 