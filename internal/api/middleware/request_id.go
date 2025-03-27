package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RequestID() fiber.Handler {
	const headerRequestID = "X-Request-ID"

	return func(c *fiber.Ctx) error {
		requestID := c.Get(headerRequestID)
		if requestID == "" {
			requestID = uuid.New().String()
			c.Set(headerRequestID, requestID)
		}
		
		c.Locals("requestID", requestID)
		
		c.Set(headerRequestID, requestID)
		
		return c.Next()
	}
}