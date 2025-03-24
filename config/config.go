package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration settings
type Config struct {
	Environment string        `mapstructure:"ENVIRONMENT"`
	Server      ServerConfig  `mapstructure:"SERVER"`
	Database    DBConfig      `mapstructure:"DATABASE"`
	JWT         JWTConfig     `mapstructure:"JWT"`
	Log         LogConfig     `mapstructure:"LOG"`
}

// ServerConfig holds server configuration settings
type ServerConfig struct {
	Port string `mapstructure:"PORT"`
}

// DBConfig holds database configuration settings
type DBConfig struct {
	Driver   string `mapstructure:"DRIVER"`
	Source   string `mapstructure:"SOURCE"`
	MigrationURL string `mapstructure:"MIGRATION_URL"`
}

// JWTConfig holds JWT configuration settings
type JWTConfig struct {
	SecretKey     string        `mapstructure:"SECRET_KEY"`
	TokenDuration time.Duration `mapstructure:"TOKEN_DURATION"`
}

// LogConfig holds logging configuration settings
type LogConfig struct {
	Level string `mapstructure:"LEVEL"`
	Path  string `mapstructure:"PATH"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}