package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string        `mapstructure:"ENVIRONMENT"`
	Server      ServerConfig  `mapstructure:"SERVER"`
	Database    DBConfig      `mapstructure:"DATABASE"`
	JWT         JWTConfig     `mapstructure:"JWT"`
	Log         LogConfig     `mapstructure:"LOG"`
}

type ServerConfig struct {
	Port string `mapstructure:"PORT"`
}

type DBConfig struct {
	Driver   string `mapstructure:"DRIVER"`
	Source   string `mapstructure:"SOURCE"`
	MigrationURL string `mapstructure:"MIGRATION_URL"`
}

type JWTConfig struct {
	SecretKey     string        `mapstructure:"SECRET_KEY"`
	TokenDuration time.Duration `mapstructure:"TOKEN_DURATION"`
	RefreshSecretKey string        `mapstructure:"REFRESH_SECRET_KEY"`
	RefreshDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

type LogConfig struct {
	Level string `mapstructure:"LEVEL"`
	Path  string `mapstructure:"PATH"`
}

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