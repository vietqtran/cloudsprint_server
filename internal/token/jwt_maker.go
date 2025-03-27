package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
	TokenType string    `json:"token_type"` 
}

func NewPayload(userID uuid.UUID, username string, duration time.Duration, tokenType string) (*Payload, error) {
	tokenID := fmt.Sprintf("%d_%s_%d", userID, username, time.Now().UnixNano())

	payload := &Payload{
		ID:        tokenID,
		UserID:    userID.String(),
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(time.Hour * duration),
		TokenType: tokenType,
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

type Maker interface {
	CreateToken(userID uuid.UUID, username string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
	CreateRefreshToken(userID uuid.UUID, username string, duration time.Duration) (string, *Payload, error)
	VerifyRefreshToken(refreshToken string) (*Payload, error)
}

type JWTMaker struct {
	accessTokenSecretKey  string
	refreshTokenSecretKey string
}

func NewJWTMaker(accessTokenSecretKey, refreshTokenSecretKey string) (Maker, error) {
	if len(accessTokenSecretKey) < 32 {
		return nil, fmt.Errorf("invalid access token key size: must be at least 32 characters")
	}
	
	if len(refreshTokenSecretKey) < 32 {
		return nil, fmt.Errorf("invalid refresh token key size: must be at least 32 characters")
	}
	
	return &JWTMaker{
		accessTokenSecretKey:  accessTokenSecretKey,
		refreshTokenSecretKey: refreshTokenSecretKey,
	}, nil
}

func (maker *JWTMaker) CreateToken(userID uuid.UUID, username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, username, duration, "access")
	if err != nil {
		return "", nil, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.accessTokenSecretKey))
	return token, payload, err
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	return maker.verifyTokenWithSecret(token, maker.accessTokenSecretKey)
}

func (maker *JWTMaker) CreateRefreshToken(userID uuid.UUID, username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, username, duration, "refresh")
	if err != nil {
		return "", nil, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.refreshTokenSecretKey))
	return token, payload, err
}

func (maker *JWTMaker) VerifyRefreshToken(refreshToken string) (*Payload, error) {
	payload, err := maker.verifyTokenWithSecret(refreshToken, maker.refreshTokenSecretKey)
	if err != nil {
		return nil, err
	}
	
	if payload.TokenType != "refresh" {
		return nil, ErrInvalidToken
	}
	
	return payload, nil
}

func (maker *JWTMaker) verifyTokenWithSecret(token string, secretKey string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		fmt.Println("Error verifying token:", err)
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}