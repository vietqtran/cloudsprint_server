package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// NewLogger creates a new logger middleware
func NewLogger(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start timer
		start := time.Now()

		// Get request IP
		ip := c.IP()

		// Process request
		err := c.Next()

		// Calculate response time
		responseTime := time.Since(start)

		// Get status code
		statusCode := c.Response().StatusCode()

		// Get route path
		path := c.Path()

		// Get request method
		method := c.Method()

		// Get request headers (optional)
		userAgent := c.Get("User-Agent")

		// Get error
		error := c.Context().Err()
		if error != nil {
			logger.Error("request failed", zap.String("error", error.Error()))
		}

		// Log request
		logger.Info("request processed",
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", ip),
			zap.String("user-agent", userAgent),
			zap.Duration("latency", responseTime),
		)

		return err
	}
}
