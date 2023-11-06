package logger

import (
	"go.uber.org/zap"
)

func NewZapLogger(cfg LogConfig) (logger *zap.Logger) {
	switch cfg.Environment {
	case "production":
		logger, _ = zap.NewProductionConfig().Build()
	default:
		logger, _ = zap.NewDevelopmentConfig().Build()
	}
	return logger
}
