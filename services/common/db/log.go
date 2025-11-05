package db

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type ormLog struct {
	LogLevel logger.LogLevel
}

func NewOrmLog() *ormLog {
	return &ormLog{
		LogLevel: logger.Info,
	}
}

func (l *ormLog) WithLogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

func (l *ormLog) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

func (l *ormLog) Info(ctx context.Context, format string, v ...interface{}) {
	if l.LogLevel < logger.Info {
		return
	}

	zap.S().Infof(format, v...)
}

func (l *ormLog) Warn(ctx context.Context, format string, v ...interface{}) {
	if l.LogLevel < logger.Warn {
		return
	}
	zap.S().Warnf(format, v...)
}

func (l *ormLog) Error(ctx context.Context, format string, v ...interface{}) {
	if l.LogLevel < logger.Error {
		return
	}
	zap.S().Errorf(format, v...)
}

func (l *ormLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	fmt.Println(elapsed, sql, rows)
	zap.S().Infof("[%.3fms] [rows:%v] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
}
