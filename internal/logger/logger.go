package logger

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new logger instance
func NewLogger(environment string) (*zap.Logger, error) {
	var config zap.Config

	if environment == "production" {
		// Production config
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		// Development config
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	}

	// Create logs directory if it doesn't exist
	logDir := "./logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			return nil, err
		}
	}

	// Configure error log output file
	today := time.Now().Format("2006-01-02")
	logFile := filepath.Join(logDir, today+".log")
	
	// Create or open log file
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Create core for console output
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config.EncoderConfig),
		zapcore.AddSync(os.Stdout),
		config.Level,
	)

	// Create core for file output
	fileEncoder := zapcore.NewJSONEncoder(config.EncoderConfig)
	fileCore := zapcore.NewCore(
		fileEncoder,
		zapcore.AddSync(file),
		zap.ErrorLevel, // Only log errors and fatals to file
	)

	// Combine cores
	core := zapcore.NewTee(consoleCore, fileCore)

	// Create logger
	logger := zap.New(core)

	return logger, nil
}