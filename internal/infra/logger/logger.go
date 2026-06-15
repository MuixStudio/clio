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

package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	global *zap.Logger
	once   sync.Once
)

type BuildOption func(*buildOptions)

type buildOptions struct {
	callerSkip int
	addCaller  bool
}

func WithCallerSkip(skip int) BuildOption {
	return func(o *buildOptions) { o.callerSkip = skip }
}

func WithCaller(enabled bool) BuildOption {
	return func(o *buildOptions) { o.addCaller = enabled }
}

// Options logger配置
type Options struct {
	Level      string // debug, info, warn, error
	Format     string // json, console
	OutputPath string // 输出文件路径，空则只输出到控制台
	MaxSize    int    // 单个日志文件最大 MB，默认 100
	MaxBackups int    // 保留旧日志文件数，默认 3
	MaxAge     int    // 保留天数，默认 7
	Compress   bool   // 是否压缩旧日志
}

func defaultOptions() Options {
	return Options{
		Level:      "info",
		Format:     "json",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   false,
	}
}

// Init 初始化全局 logger，应在 bootstrap 阶段调用一次
func Init(opts Options, extras ...BuildOption) {
	once.Do(func() {
		global = New(opts, extras...)
	})
}

// MustInit 允许重复初始化（用于测试或配置热重载）
func MustInit(opts Options, extras ...BuildOption) {
	global = New(opts, extras...)
}

func New(opts Options, extras ...BuildOption) *zap.Logger {
	bo := &buildOptions{addCaller: true, callerSkip: 1}
	for _, fn := range extras {
		fn(bo)
	}

	level := zapcore.InfoLevel
	_ = level.UnmarshalText([]byte(opts.Level))

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if opts.Format == "console" {
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	cores := []zapcore.Core{
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
	}

	if opts.OutputPath != "" {
		maxSize := opts.MaxSize
		if maxSize == 0 {
			maxSize = 100
		}
		fileWriter := &lumberjack.Logger{
			Filename:   opts.OutputPath,
			MaxSize:    maxSize,
			MaxBackups: opts.MaxBackups,
			MaxAge:     opts.MaxAge,
			Compress:   opts.Compress,
		}
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(fileWriter), level))
	}

	var zapOpts []zap.Option
	if bo.addCaller {
		zapOpts = append(zapOpts, zap.AddCaller(), zap.AddCallerSkip(bo.callerSkip))
	}
	return zap.New(zapcore.NewTee(cores...), zapOpts...)
}

// --- Global Logger ---

func L() *zap.Logger {
	if global == nil {
		// 未初始化时返回默认 logger，避免 nil panic
		global = New(defaultOptions())
	}
	return global
}

func Debug(msg string, fields ...zap.Field) { L().Debug(msg, fields...) }
func Info(msg string, fields ...zap.Field)  { L().Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { L().Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field) { L().Error(msg, fields...) }
func Fatal(msg string, fields ...zap.Field) { L().Fatal(msg, fields...) }

// With 返回带固定字段的子 logger，适合在 service/handler 层使用
func With(fields ...zap.Field) *zap.Logger {
	return L().With(fields...)
}

// Sync 在程序退出前调用，刷新缓冲区
func Sync() error {
	return L().Sync()
}
