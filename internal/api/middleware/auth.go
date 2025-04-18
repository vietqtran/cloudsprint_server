package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"cloud-sprint/config"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/internal/token"
	"cloud-sprint/pkg/util"
)

type AuthMiddleware struct {
	tokenMaker token.Maker
	tokenType  string
	tokenName  string
	store      db.Querier
	config     config.Config
}

func NewAuthMiddleware(tokenMaker token.Maker, tokenType string, store db.Querier, config config.Config) fiber.Handler {
	tokenName := "Authorization"
	if tokenType == "refresh" {
		tokenName = "Refresh"
	}

	middleware := &AuthMiddleware{
		tokenMaker: tokenMaker,
		tokenType:  tokenType,
		tokenName:  tokenName,
		store:      store,
		config:     config,
	}

	return middleware.Handle
}

func (middleware *AuthMiddleware) Handle(c *fiber.Ctx) error {
	tokenString := c.Cookies(middleware.tokenName)

	if tokenString == "" {
		authHeader := c.Get(middleware.tokenName)
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			} else if len(parts) == 1 {
				tokenString = parts[0]
			}
		}
	}

	if tokenString == "" {
		return fiber.NewError(fiber.StatusUnauthorized, middleware.tokenType+" token is missing")
	}

	var payload *token.Payload
	var err error

	if middleware.tokenType == "refresh" {
		payload, err = middleware.tokenMaker.VerifyRefreshToken(tokenString)
	} else {
		payload, err = middleware.tokenMaker.VerifyToken(tokenString)
		
		if err == token.ErrExpiredToken && middleware.store != nil {
			refreshToken := c.Cookies("Refresh")
			if refreshToken == "" {
				refreshHeader := c.Get("Refresh")
				if refreshHeader != "" {
					parts := strings.Split(refreshHeader, " ")
					if len(parts) == 2 && parts[0] == "Bearer" {
						refreshToken = parts[1]
					} else if len(parts) == 1 {
						refreshToken = parts[0]
					}
				}
			}
			
			if refreshToken != "" {
				refreshPayload, verifyErr := middleware.tokenMaker.VerifyRefreshToken(refreshToken)
				if verifyErr == nil {
					userID, parseErr := uuid.Parse(refreshPayload.UserID)
					if parseErr == nil {
						newAccessToken, _, tokenErr := middleware.tokenMaker.CreateToken(
							userID,
							refreshPayload.Email,
							middleware.config.JWT.TokenDuration,
						)
						
						if tokenErr == nil {
							newRefreshToken, refreshPayload, tokenErr := middleware.tokenMaker.CreateRefreshToken(
								userID,
								refreshPayload.Email,
								middleware.config.JWT.RefreshDuration,
							)
							
							if tokenErr == nil {
								util.SetHttpOnlyCookie(c, util.SetCookieData{
									Name:      "Authorization",
									Token:     newAccessToken,
									ExpiresAt: int(middleware.config.JWT.TokenDuration.Seconds()),
									ENV:       middleware.config.Environment,
								})
								
								util.SetHttpOnlyCookie(c, util.SetCookieData{
									Name:      "Refresh",
									Token:     newRefreshToken,
									ExpiresAt: int(middleware.config.JWT.RefreshDuration.Seconds()),
									ENV:       middleware.config.Environment,
								})
								
								c.Set("Authorization", "Bearer "+newAccessToken)
								
								c.Locals("current_user_id", refreshPayload.UserID)
								c.Locals("current_user_email", refreshPayload.Email)
								
								return c.Next()
							}
						}
					}
				}
			}
			
			return fiber.NewError(fiber.StatusUnauthorized, middleware.tokenType+" token has expired")
		}
	}

	if err != nil {
		if err == token.ErrExpiredToken {
			return fiber.NewError(fiber.StatusUnauthorized, middleware.tokenType+" token has expired")
		}
		return fiber.NewError(fiber.StatusUnauthorized, "invalid "+middleware.tokenType+" token")
	}

	c.Locals("current_user_id", payload.UserID)
	c.Locals("current_user_email", payload.Email)
	if middleware.tokenType == "refresh" {
		c.Locals("refresh_token", tokenString)
	}

	return c.Next()
}
