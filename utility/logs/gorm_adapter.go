package logs

import (
	"context"
	"github.com/apsdehal/go-logger"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

type GormLogger struct {
	logger *logger.Logger
	level  gormLogger.LogLevel
}

func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return &GormLogger{
		logger: l.logger,
		level:  level,
	}
}

func (l *GormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	if l.level >= gormLogger.Info {
		l.logger.Infof(s, i...)
	}
}

func (l *GormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	if l.level >= gormLogger.Warn {
		l.logger.Warningf(s, i...)
	}
}

func (l *GormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	if l.level >= gormLogger.Error {
		l.logger.Errorf(s, i...)
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.level <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rowsAffected := fc()

	if err != nil {
		l.logger.Errorf("gorm: query: %s, exectime: %s, error: %s", sql, elapsed, err)
	} else {
		l.logger.Debugf("gorm: query: %s, exectime: %s, row affected: %d", sql, elapsed, rowsAffected)
	}
}
