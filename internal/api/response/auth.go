package response

import (
	"time"

	db "cloud-sprint/internal/db/sqlc"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID     `json:"id"`
	Username  string    		`json:"username"`
	Email     string    		`json:"email"`
	FullName  string    		`json:"full_name"`
	CreatedAt time.Time 		`json:"created_at"`
	UpdatedAt time.Time 		`json:"updated_at"`
}

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

type LoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	SessionID    string       `json:"session_id"`
}

func NewLoginResponse(user db.User, accessToken, refreshToken, sessionID string) LoginResponse {
	return LoginResponse{
		User:         NewUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		SessionID:    sessionID,
	}
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
