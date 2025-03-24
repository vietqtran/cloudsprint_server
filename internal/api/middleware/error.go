package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// NewErrorHandler creates a new error handling middleware
func NewErrorHandler(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Continue to the next middleware/handler
		err := c.Next()

		// If there's no error, return nil
		if err == nil {
			return nil
		}

		// Get error code
		code := fiber.StatusInternalServerError

		// Check if it's a Fiber error
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		// Log the error
		logger.Error("request error",
			zap.Int("status", code),
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
			zap.String("ip", c.IP()),
			zap.Error(err),
		)

		// Return JSON error response
		return c.Status(code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
}
