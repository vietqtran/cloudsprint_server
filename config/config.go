package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment     string
	Server          ServerConfig
	Database        DBConfig
	JWT             JWTConfig
	Log             LogConfig
	Email           EmailConfig
	FrontendBaseURL string
	OAuth           OAuthConfig
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	Driver       string
	Source       string
	MigrationURL string
}

type JWTConfig struct {
	SecretKey        string
	TokenDuration    time.Duration
	RefreshSecretKey string
	RefreshDuration  time.Duration
}

type LogConfig struct {
	Level string
	Path  string
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
	TemplatesDir string
}

type OAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	GitHubClientID     string
	GitHubClientSecret string
	GitHubRedirectURL  string
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, fmt.Errorf("error loading .env file: %w", err)
	}

	tokenDuration, err := parseDuration("JWT_TOKEN_DURATION")
	if err != nil {
		return Config{}, err
	}

	refreshDuration, err := parseDuration("JWT_REFRESH_TOKEN_DURATION")
	if err != nil {
		return Config{}, err
	}

	smtpPort, err := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	if err != nil {
		smtpPort = 587
	}

	config := Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DBConfig{
			Driver:       getEnv("DB_DRIVER", "postgres"),
			Source:       getEnv("DB_SOURCE", ""),
			MigrationURL: getEnv("DB_MIGRATION_URL", "file://db/migration"),
		},
		JWT: JWTConfig{
			SecretKey:        getEnv("JWT_SECRET_KEY", ""),
			TokenDuration:    tokenDuration,
			RefreshSecretKey: getEnv("JWT_REFRESH_SECRET_KEY", ""),
			RefreshDuration:  refreshDuration,
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			Path:  getEnv("LOG_PATH", "logs"),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     smtpPort,
			SMTPUsername: getEnv("SMTP_USERNAME", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			FromEmail:    getEnv("EMAIL_FROM", ""),
			FromName:     getEnv("EMAIL_FROM_NAME", ""),
			TemplatesDir: getEnv("EMAIL_TEMPLATES_DIR", "./templates/emails"),
		},
		FrontendBaseURL: getEnv("FRONTEND_BASE_URL", "http://localhost:3000"),
		OAuth: OAuthConfig{
			GoogleClientID:     getEnv("CLIENT_ID", ""),
			GoogleClientSecret: getEnv("CLIENT_SECRET", ""),
			GoogleRedirectURL:  getEnv("REDIRECT_URL", ""),
			GitHubClientID:     getEnv("GitHubClientID", ""),
			GitHubClientSecret: getEnv("GitHubClientSecret", ""),
			GitHubRedirectURL:  getEnv("GitHubRedirectURL", ""),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func parseDuration(key string) (time.Duration, error) {
	durationStr := os.Getenv(key)

	if hours, err := strconv.Atoi(durationStr); err == nil {
		return time.Duration(hours) * time.Hour, nil
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, fmt.Errorf("invalid duration for %s: %w", key, err)
	}

	return duration, nil
}
