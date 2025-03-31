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
	accessToken := c.Cookies(middleware.cookieName)
	if accessToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "authentication cookie is missing")
	}
	payload, err := middleware.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		if err == token.ErrExpiredToken {
			return fiber.NewError(fiber.StatusUnauthorized, "token has expired")
		}
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}
	c.Locals("current_user_id", payload.UserID)
	c.Locals("current_user_email", payload.Email)

	return c.Next()
}
