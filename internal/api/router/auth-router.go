package router

import (
	"github.com/gofiber/fiber/v2"

	"cloud-sprint/config"
	"cloud-sprint/internal/api/handler"
	"cloud-sprint/internal/api/middleware"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/internal/service"
	"cloud-sprint/internal/token"
)

func SetupAuthRoutes(api fiber.Router, store db.Querier, tokenMaker token.Maker, config config.Config) {
	emailService := service.NewEmailService(config.Email)
	authHandler := handler.NewAuthHandler(store, tokenMaker, config, emailService)

	authMiddleware := middleware.NewAuthMiddleware(tokenMaker, "access")
	refreshMiddleware := middleware.NewAuthMiddleware(tokenMaker, "refresh")

	auth := api.Group("/auth")
	auth.Post("/sign-up", authHandler.SignUp)
	auth.Post("/sign-in", authHandler.SignIn)
	auth.Post("/refresh", refreshMiddleware, authHandler.RefreshToken)
	auth.Get("/me", authMiddleware, authHandler.Me)

	passwordHandler := handler.NewPasswordHandler(store, tokenMaker, config, emailService)
	auth.Post("/forgot-password", passwordHandler.ForgotPassword)
	auth.Post("/verify-reset-token", passwordHandler.VerifyResetToken)
	auth.Post("/reset-password", passwordHandler.ResetPassword)

	emailVerificationHandler := handler.NewEmailVerificationHandler(store, config, tokenMaker, emailService)
	verifyEmail := auth.Group("/verify-email")
	verifyEmail.Post("/send-otp", emailVerificationHandler.SendOTP)
	verifyEmail.Post("/verify", emailVerificationHandler.VerifyOTP)
	verifyEmail.Get("/status", authMiddleware, emailVerificationHandler.CheckVerificationStatus)
}
