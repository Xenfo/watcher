package logger

import (
	"os"

	"go.uber.org/zap"
)

func Create() (*zap.Logger, error) {
	var logger *zap.Logger

	if os.Getenv("WATCHER_PRODUCTION") == "" {
		devLogger, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}

		logger = devLogger
	} else {
		prodLogger, err := zap.NewProduction()
		if err != nil {
			return nil, err
		}

		logger = prodLogger
	}

	return logger, nil
}
