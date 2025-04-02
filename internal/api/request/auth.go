package request

import (
	"errors"
	"net/mail"
	"strings"
)

type SignUpRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (r *SignUpRequest) Validate() error {
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

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *SignInRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}

	if r.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

type RefreshTokenRequest struct {
	SessionID string `json:"sessionId"`
}

func (r *RefreshTokenRequest) Validate() error {
	if r.SessionID == "" {
		return errors.New("session ID is required")
	}

	return nil
}

type UpdateUserRequest struct {
	Email     *string `json:"email,omitempty"`
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
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
