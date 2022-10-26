package utils

import (
	"os"
	"sync"

	"github.com/google/wire"
	"go.uber.org/zap"
)

var (
	loggerSugar *zap.SugaredLogger
	loggerInit  sync.Once
)

var SuperSet = wire.NewSet(
	ProvideLogger,
)

func injectLogger() *zap.SugaredLogger {
	var logger *zap.Logger
	var err error
	if os.Getenv("ENV") == "LOCAL" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}
	return logger.Sugar()
}

func ProvideLogger() *zap.SugaredLogger {
	loggerInit.Do(func() {
		loggerSugar = injectLogger()
	})
	return loggerSugar
}
