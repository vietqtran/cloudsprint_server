package response

import (
	"time"

	db "cloud-sprint/internal/db/sqlc"
)

// UserResponse represents a user response
type UserResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUserResponse creates a new user response from a user
func NewUserResponse(user db.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// LoginResponse represents a login response
type LoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	SessionID    string       `json:"session_id"`
}

// NewLoginResponse creates a new login response
func NewLoginResponse(user db.User, accessToken, refreshToken, sessionID string) LoginResponse {
	return LoginResponse{
		User:         NewUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		SessionID:    sessionID,
	}
}

// RefreshTokenResponse represents a refresh token response
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}
