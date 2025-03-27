package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func NewErrorHandler(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err == nil {
			return nil
		}

		code := fiber.StatusInternalServerError

		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		logger.Error("request error",
			zap.Int("status", code),
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
			zap.String("ip", c.IP()),
			zap.Error(err),
		)

		return c.Status(code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
}
