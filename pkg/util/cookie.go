package util

import "github.com/gofiber/fiber/v2"

type SetCookieData struct {
	Name      string
	Token     string
	ExpiresAt int
	ENV       string
}

func SetHttpOnlyCookie(c *fiber.Ctx, data SetCookieData) {
	cookie := &fiber.Cookie{
		Name:     data.Name,
		Value:    data.Token,
		Path:     "/",
		MaxAge:   data.ExpiresAt,
		HTTPOnly: true,
		Secure:   data.ENV == "production",
		SameSite: "Strict",
	}

	c.Cookie(cookie)
}
