package request

import (
	"errors"
	"net/mail"
)

type SendEmailOTPRequest struct {
	Email string `json:"email"`
}

func (r *SendEmailOTPRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}

	if _, err := mail.ParseAddress(r.Email); err != nil {
		return errors.New("invalid email address")
	}

	return nil
}

type VerifyEmailOTPRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func (r *VerifyEmailOTPRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}

	if _, err := mail.ParseAddress(r.Email); err != nil {
		return errors.New("invalid email address")
	}

	if r.OTP == "" {
		return errors.New("OTP is required")
	}

	return nil
}
