package request

import (
	"errors"
	"net/mail"
)

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func (r *ForgotPasswordRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}

	if _, err := mail.ParseAddress(r.Email); err != nil {
		return errors.New("invalid email address")
	}

	return nil
}

type VerifyResetTokenRequest struct {
	Token string `json:"token"`
}

func (r *VerifyResetTokenRequest) Validate() error {
	if r.Token == "" {
		return errors.New("token is required")
	}

	return nil
}

type ResetPasswordRequest struct {
	Token           string `json:"token"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func (r *ResetPasswordRequest) Validate() error {
	if r.Token == "" {
		return errors.New("token is required")
	}

	if r.Password == "" {
		return errors.New("password is required")
	}

	if r.ConfirmPassword == "" {
		return errors.New("confirm password is required")
	}

	if r.Password != r.ConfirmPassword {
		return errors.New("passwords do not match")
	}

	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	return nil
}
