package logs

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	gormLogger "gorm.io/gorm/logger"
	"os"
	"strings"
)

var defaultLogger *Logger

type Config struct {
	LogLevel string
}

type Logger struct {
	*logger.Logger
	level logger.LogLevel
}

func Init(cfg Config) {
	log, err := logger.New("default", 1, os.Stdout)
	if err != nil {
		panic(fmt.Errorf("Cannot create default logger: %v\n", err))
	}
	defaultLogger = &Logger{Logger: log, level: toLogLevel(cfg.LogLevel)}
	defaultLogger.SetLogLevel(defaultLogger.level)
	defaultLogger.SetFormat("%{time} %{file}:%{line} [%{level}] ▶ %{message}")
	defaultLogger.Debug("Default logger initialized")
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

func Log() *Logger {
	return defaultLogger
}

func (l *Logger) ToGormLogger() *GormLogger {
	var level gormLogger.LogLevel
	switch l.level {
	case logger.DebugLevel, logger.InfoLevel:
		level = gormLogger.Info
	case logger.NoticeLevel, logger.WarningLevel:
		level = gormLogger.Warn
	case logger.ErrorLevel:
		level = gormLogger.Error
	case logger.CriticalLevel:
		level = gormLogger.Silent
	}
	return &GormLogger{logger: l.Logger, level: level}
}

func Debug(s string) {
	defaultLogger.Debug(s)
}

func Info(s string) {
	defaultLogger.Info(s)
}

func Warn(s string) {
	defaultLogger.Warning(s)
}

func Error(s string) {
	defaultLogger.Error(s)
}

func Fatal(s string) {
	defaultLogger.Fatal(s)
}

func Debugf(s string, args ...any) {
	defaultLogger.Debugf(s, args...)
}

func Infof(s string, args ...any) {
	defaultLogger.Infof(s, args...)
}

func Warnf(s string, args ...any) {
	defaultLogger.Warningf(s, args...)
}

func Errorf(s string, args ...any) {
	defaultLogger.Errorf(s, args...)
}

func Fatalf(s string, args ...any) {
	defaultLogger.Fatalf(s, args...)
}
