package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"cloud-sprint/internal/token"
)

type AuthMiddleware struct {
	tokenMaker token.Maker
	tokenType  string
	tokenName  string
}

func NewAuthMiddleware(tokenMaker token.Maker, tokenType string) fiber.Handler {
	tokenName := "Authorization"
	if tokenType == "refresh" {
		tokenName = "Refresh"
	}

	middleware := &AuthMiddleware{
		tokenMaker: tokenMaker,
		tokenType:  tokenType,
		tokenName:  tokenName,
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
