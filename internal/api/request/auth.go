package request

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

// Validate validates the register request
func (r *RegisterRequest) Validate() error {
	// Validate username
	if len(r.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}

	if strings.Contains(r.Username, " ") {
		return errors.New("username cannot contain spaces")
	}

	// Validate email
	if _, err := mail.ParseAddress(r.Email); err != nil {
		return errors.New("invalid email address")
	}

	// Validate password
	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	// Validate full name
	if len(r.FullName) < 2 {
		return errors.New("full name must be at least 2 characters long")
	}

	return nil
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate validates the login request
func (r *LoginRequest) Validate() error {
	// Validate username
	if r.Username == "" {
		return errors.New("username is required")
	}

	// Validate password
	if r.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	SessionID    string `json:"session_id"`
}

// Validate validates the refresh token request
func (r *RefreshTokenRequest) Validate() error {
	// Validate refresh token
	if r.RefreshToken == "" {
		return errors.New("refresh token is required")
	}

	// Validate session ID
	if r.SessionID == "" {
		return errors.New("session ID is required")
	}

	return nil
}

// UpdateUserRequest represents an update user request
type UpdateUserRequest struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
	FullName *string `json:"full_name,omitempty"`
}

// Validate validates the update user request
func (r *UpdateUserRequest) Validate() error {
	// Check if at least one field is provided
	if r.Username == nil && r.Email == nil && r.Password == nil && r.FullName == nil {
		return errors.New("at least one field must be provided")
	}

	// Validate username if provided
	if r.Username != nil {
		if len(*r.Username) < 3 {
			return errors.New("username must be at least 3 characters long")
		}

		if strings.Contains(*r.Username, " ") {
			return errors.New("username cannot contain spaces")
		}
	}

	// Validate email if provided
	if r.Email != nil {
		if _, err := mail.ParseAddress(*r.Email); err != nil {
			return errors.New("invalid email address")
		}
	}

	// Validate password if provided
	if r.Password != nil {
		if len(*r.Password) < 6 {
			return fmt.Errorf("password must be at least 6 characters long")
		}
	}

	// Validate full name if provided
	if r.FullName != nil {
		if len(*r.FullName) < 2 {
			return errors.New("full name must be at least 2 characters long")
		}
	}

	return nil
}
