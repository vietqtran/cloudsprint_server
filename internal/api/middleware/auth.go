package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"cloud-sprint/internal/token"
)

type AuthMiddleware struct {
	tokenMaker token.Maker
}

func NewAuthMiddleware(tokenMaker token.Maker) fiber.Handler {
	middleware := &AuthMiddleware{
		tokenMaker: tokenMaker,
	}

	return middleware.Handle
}

func (middleware *AuthMiddleware) Handle(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if len(authHeader) == 0 {
		return fiber.NewError(fiber.StatusUnauthorized, "authorization header is required")
	}

	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != "bearer" {
		return fiber.NewError(fiber.StatusUnauthorized, "unsupported authorization type")
	}

	accessToken := fields[1]
	payload, err := middleware.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		if err == token.ErrExpiredToken {
			return fiber.NewError(fiber.StatusUnauthorized, "token has expired")
		}
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}

	c.Locals("userID", payload.UserID)
	c.Locals("username", payload.Username)

	return c.Next()
}
