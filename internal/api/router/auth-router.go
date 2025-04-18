package router

import (
	"github.com/gofiber/fiber/v2"

	"cloud-sprint/config"
	"cloud-sprint/internal/api/handler"
	db "cloud-sprint/internal/db/sqlc"
	"cloud-sprint/internal/service"
	"cloud-sprint/internal/token"
)

func SetupAuthRoutes(api fiber.Router, store db.Querier, tokenMaker token.Maker, config config.Config, authMiddleware fiber.Handler, refreshMiddleware fiber.Handler) {
	emailService := service.NewEmailService(config.Email)
	googleService := service.NewGoogleService(config)
	githubService := service.NewGitHubService(config)

	authHandler := handler.NewAuthHandler(store, tokenMaker, config, emailService)

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
	verifyEmail.Get("/status", emailVerificationHandler.CheckVerificationStatus)

	googleAuthHandler := handler.NewGoogleAuthHandler(store, tokenMaker, config, emailService, googleService)
	googleAuth := auth.Group("/google")
	googleAuth.Get("/auth", googleAuthHandler.GoogleAuth)
	googleAuth.Get("/callback", googleAuthHandler.GoogleCallback)

	githubAuthHandler := handler.NewGitHubAuthHandler(store, tokenMaker, config, emailService, githubService)
	githubAuth := auth.Group("/github")
	githubAuth.Get("/auth", githubAuthHandler.GitHubAuth)
	githubAuth.Get("/callback", githubAuthHandler.GitHubCallback)
}
