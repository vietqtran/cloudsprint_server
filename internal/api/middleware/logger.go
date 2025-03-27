package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func NewLogger(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		ip := c.IP()

		err := c.Next()

		responseTime := time.Since(start)

		statusCode := c.Response().StatusCode()

		path := c.Path()

		method := c.Method()

		userAgent := c.Get("User-Agent")

		error := c.Context().Err()
		if error != nil {
			logger.Error("request failed", zap.String("error", error.Error()))
		}

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
