package logs

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"os"
)

var defaultLogger *logger.Logger

func Init() {
	log, err := logger.New("default", 1, os.Stdout)
	if err != nil {
		panic(fmt.Errorf("Cannot create default logger: %v\n", err))
	}
	defaultLogger = log
	defaultLogger.SetLogLevel(logger.DebugLevel)
	defaultLogger.SetFormat("%{time} %{file}:%{line} [%{level}] ▶ %{message}")
	defaultLogger.Debug("Default logger initialized")
}

func Log() *logger.Logger {
	return defaultLogger
}
