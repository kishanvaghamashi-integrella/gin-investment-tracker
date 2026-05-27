package util

import (
	"os"

	"go.uber.org/zap"
)

var (
	baseLogger = zap.NewNop()
	Logger     = baseLogger.Sugar()
)

func InitLogger() {
	var cfg zap.Config

	if os.Getenv("APP_ENV") == "production" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	cfg.DisableStacktrace = true

	logger, err := cfg.Build()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	baseLogger = logger
	Logger = logger.Sugar()
}

func SyncLogger() error {
	if baseLogger == nil {
		return nil
	}

	return baseLogger.Sync()
}
