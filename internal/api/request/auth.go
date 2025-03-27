package request

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

// Validate validates the register request
func (r *RegisterRequest) Validate() error {
	if len(r.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}

	if strings.Contains(r.Username, " ") {
		return errors.New("username cannot contain spaces")
	}

	if _, err := mail.ParseAddress(r.Email); err != nil {
		return errors.New("invalid email address")
	}

	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	if len(r.FullName) < 2 {
		return errors.New("full name must be at least 2 characters long")
	}

	return nil
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if r.Username == "" {
		return errors.New("username is required")
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
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
	FullName *string `json:"full_name,omitempty"`
}

func (r *UpdateUserRequest) Validate() error {
	if r.Username == nil && r.Email == nil && r.Password == nil && r.FullName == nil {
		return errors.New("at least one field must be provided")
	}

	if r.Username != nil {
		if len(*r.Username) < 3 {
			return errors.New("username must be at least 3 characters long")
		}

		if strings.Contains(*r.Username, " ") {
			return errors.New("username cannot contain spaces")
		}
	}

	if r.Email != nil {
		if _, err := mail.ParseAddress(*r.Email); err != nil {
			return errors.New("invalid email address")
		}
	}

	if r.Password != nil {
		if len(*r.Password) < 6 {
			return fmt.Errorf("password must be at least 6 characters long")
		}
	}

	if r.FullName != nil {
		if len(*r.FullName) < 2 {
			return errors.New("full name must be at least 2 characters long")
		}
	}

	return nil
}
