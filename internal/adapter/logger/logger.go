package logger

import "go.uber.org/zap"

func NewLogger(appName string) *zap.Logger {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {

		}
	}(logger)

	return logger.Named(appName)
}
