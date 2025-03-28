package request

import (
	"errors"
	"net/mail"
	"strings"
)

type RegisterRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
}

// Validate validates the register request
func (r *RegisterRequest) Validate() error {
	if strings.Contains(r.FirstName, " ") {
		return errors.New("first name cannot contain spaces")
	}

	if strings.Contains(r.LastName, " ") {
		return errors.New("last name cannot contain spaces")
	}

	if _, err := mail.ParseAddress(r.Email); err != nil {
		return errors.New("invalid email address")
	}

	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}

	if r.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	SessionID    string `json:"session_id"`
}

func (r *RefreshTokenRequest) Validate() error {
	if r.RefreshToken == "" {
		return errors.New("refresh token is required")
	}

	if r.SessionID == "" {
		return errors.New("session ID is required")
	}

	return nil
}

type UpdateUserRequest struct {
	Email      *string `json:"email,omitempty"`
	FirstName  *string `json:"first_name,omitempty"`
	LastName   *string `json:"last_name,omitempty"`
}

func (r *UpdateUserRequest) Validate() error {
	if strings.Contains(*r.FirstName, " ") {
		return errors.New("first name cannot contain spaces")
	}

	if strings.Contains(*r.LastName, " ") {
		return errors.New("last name cannot contain spaces")
	}

	if r.Email != nil {
		if _, err := mail.ParseAddress(*r.Email); err != nil {
			return errors.New("invalid email address")
		}
	}

	return nil
}
