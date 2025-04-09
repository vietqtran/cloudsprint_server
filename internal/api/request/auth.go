package request

import (
	"cloud-sprint/internal/api/response"
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
		return response.BadRequest(nil, "first name cannot contain spaces", nil, nil)
	}

	if strings.Contains(r.LastName, " ") {
		return response.BadRequest(nil, "last name cannot contain spaces", nil, nil)
	}

	if _, err := mail.ParseAddress(r.Email); err != nil {
		return response.BadRequest(nil, "invalid email address", nil, nil)
	}

	if len(r.Password) < 6 {
		return response.BadRequest(nil, "password must be at least 6 characters long", nil, nil)
	}

	return nil
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *SignInRequest) Validate() error {
	if r.Email == "" {
		return response.BadRequest(nil, "email is required", nil, nil)
	}

	if r.Password == "" {
		return response.BadRequest(nil, "password is required", nil, nil)
	}

	return nil
}

type RefreshTokenRequest struct {
	SessionID string `json:"sessionId,omitempty"`
}

func (r *RefreshTokenRequest) Validate() error {
	if r.SessionID == "" {
		return response.BadRequest(nil, "Session ID is required", nil, nil)
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
