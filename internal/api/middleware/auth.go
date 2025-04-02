package middleware

import (
	"github.com/gofiber/fiber/v2"

	"cloud-sprint/internal/token"
)

type AuthMiddleware struct {
	tokenMaker token.Maker
	cookieName string
}

func NewAuthMiddleware(tokenMaker token.Maker, cookieName string) fiber.Handler {
	if cookieName == "" {
		cookieName = "Authorization"
	}

	middleware := &AuthMiddleware{
		tokenMaker: tokenMaker,
		cookieName: cookieName,
	}

	return middleware.Handle
}

func (middleware *AuthMiddleware) Handle(c *fiber.Ctx) error {
	tokenString := c.Cookies(middleware.cookieName)
	if tokenString == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "token is missing")
	}
	payload, err := middleware.tokenMaker.VerifyToken(tokenString)
	if err != nil {
		if err == token.ErrExpiredToken {
			return fiber.NewError(fiber.StatusUnauthorized, "token has expired")
		}
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}
	if middleware.cookieName == "Authorization" {
		c.Locals("current_user_id", payload.UserID)
		c.Locals("current_user_email", payload.Email)
	} else {
		c.Locals("refresh_token", tokenString)
	}

	return c.Next()
}
