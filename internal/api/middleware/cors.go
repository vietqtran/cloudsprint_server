package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORS configures Cross-Origin Resource Sharing for the application
func CORS() fiber.Handler {
	return cors.New(cors.Config{
		// Instead of allowing all origins with "*", we'll explicitly handle it in code
		// to allow credentials while maintaining security
		AllowOriginsFunc: func(origin string) bool {
			// In a production environment, you'd check against a whitelist
			// For development, we'll allow all origins
			return true
		},
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",                                      // Supported methods
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With,X-Request-ID", // Allowed headers
		ExposeHeaders:    "Content-Length,Content-Type,X-Request-ID",                               // Exposed headers
		AllowCredentials: true,                                                                     // Allow credentials (cookies, etc.)
		MaxAge:           86400,                                                                    // Preflight requests are cached for 24 hours
	})
}

// CORSSimple provides a simpler CORS configuration without credentials
// Use this if you don't need to support cookies or authentication headers
func CORSSimple() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "https://localhost:8080",                                                 // Allow all origins
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",                                      // Supported methods
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With,X-Request-ID", // Allowed headers
		ExposeHeaders:    "Content-Length,Content-Type,X-Request-ID",                               // Exposed headers
		AllowCredentials: false,                                                                    // No credentials
		MaxAge:           86400,                                                                    // Preflight requests are cached for 24 hours
	})
}
