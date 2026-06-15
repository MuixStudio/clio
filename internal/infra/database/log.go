/*
 * Copyright 2026 MuixStudio
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	logger2 "github.com/muixstudio/clio/internal/infra/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type cLogger struct {
	logger        *zap.Logger
	logLevel      gormlogger.LogLevel
	slowThreshold time.Duration
}

func (c *cLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	clone := *c
	clone.logLevel = level
	return &clone
}

func (c *cLogger) Info(_ context.Context, msg string, args ...interface{}) {
	if c.logLevel < gormlogger.Info {
		return
	}
	c.logger.Info(fmt.Sprintf(msg, args...))
}

func (c *cLogger) Warn(_ context.Context, msg string, args ...interface{}) {
	if c.logLevel < gormlogger.Warn {
		return
	}
	c.logger.Warn(fmt.Sprintf(msg, args...))
}

func (c *cLogger) Error(_ context.Context, msg string, args ...interface{}) {
	if c.logLevel < gormlogger.Error {
		return
	}
	c.logger.Error(fmt.Sprintf(msg, args...))
}

func (c *cLogger) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if c.logLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("elapsed", elapsed),
	}

	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound):
		c.logger.Error("gorm query error", append(fields, zap.Error(err))...)
	case c.slowThreshold > 0 && elapsed > c.slowThreshold:
		c.logger.Warn("slow query", fields...)
	case c.logLevel >= gormlogger.Info:
		c.logger.Info("query", fields...)
	}
}

func NewLog(cfg DBConfig) gormlogger.Interface {
	l := logger2.New(
		logger2.Options{
			Level:      cfg.Log.Level,
			Format:     cfg.Log.Format,
			OutputPath: cfg.Log.OutputPath,
			MaxSize:    cfg.Log.MaxSize,
			MaxBackups: cfg.Log.MaxBackups,
			MaxAge:     cfg.Log.MaxAge,
			Compress:   cfg.Log.Compress,
		},
		logger2.WithCallerSkip(3),
	)
	return &cLogger{
		logger:        l,
		logLevel:      toGormLogLevel(cfg.Log.Level),
		slowThreshold: cfg.SlowThreshold,
	}
}

func toGormLogLevel(level string) gormlogger.LogLevel {
	switch level {
	case "silent":
		return gormlogger.Silent
	case "error":
		return gormlogger.Error
	case "warn":
		return gormlogger.Warn
	default:
		return gormlogger.Info
	}
}
