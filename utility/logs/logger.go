package logs

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"os"
	"strings"
)

type Config struct {
	LogLevel string
}

type Logger struct {
	*logger.Logger
	level logger.LogLevel
}

var TestConfig = Config{
	LogLevel: "debug",
}

func New(cfg Config) *Logger {
	log, err := logger.New("default", 1, os.Stdout)
	if err != nil {
		panic(fmt.Errorf("Cannot create default logger: %v\n", err))
	}
	defaultLogger := &Logger{Logger: log, level: toLogLevel(cfg.LogLevel)}
	defaultLogger.SetLogLevel(defaultLogger.level)
	defaultLogger.SetFormat("%{time} %{file}:%{line} [%{level}] ▶ %{message}")
	return defaultLogger
}

func toLogLevel(logLevel string) logger.LogLevel {
	logLevel = strings.ToLower(logLevel)
	switch logLevel {
	case "", "debug":
		return logger.DebugLevel
	case "info":
		return logger.InfoLevel
	case "note", "notice":
		return logger.NoticeLevel
	case "warn", "warning":
		return logger.WarningLevel
	case "err", "error":
		return logger.ErrorLevel
	case "fatal", "critical":
		return logger.CriticalLevel
	default:
		panic("invalid logLevel: " + logLevel)
	}
}
