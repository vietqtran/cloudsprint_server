package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",                                      
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With,X-Request-ID", 
		ExposeHeaders:    "Content-Length,Content-Type,X-Request-ID",                               
		AllowCredentials: true,                                                                     
		MaxAge:           86400,                                                                    
	})
}

func CORSSimple() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "*",                                                 
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",                                      
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With,X-Request-ID", 
		ExposeHeaders:    "Content-Length,Content-Type,X-Request-ID",                               
		AllowCredentials: false,                                                                    
		MaxAge:           86400,                                                                    
	})
}
