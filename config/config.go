package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string        
	Server      ServerConfig  
	Database    DBConfig      
	JWT         JWTConfig     
	Log         LogConfig     
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