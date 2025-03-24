package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestID adds a unique request ID to each request
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if request already has an ID
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			// Generate a new UUID
			requestID = uuid.New().String()
			// Add it to the request header
			c.Set("X-Request-ID", requestID)
		}
		
		// Store in locals for easy access
		c.Locals("requestID", requestID)
		
		// Add response header with the same request ID
		c.Set("X-Request-ID", requestID)
		
		return c.Next()
	}
}