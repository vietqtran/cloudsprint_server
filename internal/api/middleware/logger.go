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

		statusCode := c.Response().StatusCode()

		responseTime := time.Since(start)
		path := c.Path()
		method := c.Method()
		userAgent := c.Get("User-Agent")

		logger.Info("request processed",
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", ip),
			zap.String("user-agent", userAgent),
			zap.Duration("latency", responseTime),
			zap.Binary("request", c.Body()),
		)

		if err != nil {
			logger.Error("request error",
				zap.Error(err),
				zap.String("path", path),
				zap.String("method", method),
			)
		}

		return err
	}
}
