package logging

import (
	"fmt"
	"os"
	"path/filepath"

	"goapp/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger interface defines logging methods for both structured and unstructured logging
type Logger interface {
	// Structured logging with key-value pairs
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field) 
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	
	// Unstructured logging (printf-style)
	Infof(template string, args ...interface{})
	Debugf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	
	// Utility methods
	Sync() error
	With(fields ...zap.Field) Logger
}

// logger implements the Logger interface using zap
type logger struct {
	zap *zap.Logger
}

// New creates a new logger instance
func New(cfg config.LoggerConfig) (Logger, error) {
	l := &logger{}
	if err := l.initialize(cfg); err != nil {
		return nil, err
	}
	return l, nil
}

// initialize sets up the zap logger with the given configuration
func (l *logger) initialize(cfg config.LoggerConfig) error {
	// Create log directories if they don't exist
	if cfg.AppLogPath != "" {
		if err := os.MkdirAll(filepath.Dir(cfg.AppLogPath), 0755); err != nil {
			return fmt.Errorf("failed to create app log directory: %w", err)
		}
	}
	if cfg.ErrLogPath != "" {
		if err := os.MkdirAll(filepath.Dir(cfg.ErrLogPath), 0755); err != nil {
			return fmt.Errorf("failed to create error log directory: %w", err)
		}
	}

	// Configure encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Build cores for different outputs
	var cores []zapcore.Core

	// App log file (all levels)
	if cfg.AppLogPath != "" {
		appLogWriter := &lumberjack.Logger{
			Filename:   cfg.AppLogPath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(appLogWriter),
			zap.DebugLevel,
		))
	}

	// Error log file (error and fatal only)
	if cfg.ErrLogPath != "" {
		errorLogWriter := &lumberjack.Logger{
			Filename:   cfg.ErrLogPath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(errorLogWriter),
			zap.ErrorLevel,
		))
	}

	// Console output
	if cfg.WriteStdout {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		cores = append(cores, zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			zap.DebugLevel,
		))
	}

	// Create logger
	core := zapcore.NewTee(cores...)
	opts := []zap.Option{zap.AddCaller(), zap.AddCallerSkip(1)}
	if cfg.EnableStackTrace {
		opts = append(opts, zap.AddStacktrace(zap.ErrorLevel))
	}
	
	l.zap = zap.New(core, opts...)
	return nil
}

// Structured logging methods
func (l *logger) Info(msg string, fields ...zap.Field) {
	l.zap.Info(msg, fields...)
}

func (l *logger) Debug(msg string, fields ...zap.Field) {
	l.zap.Debug(msg, fields...)
}

func (l *logger) Error(msg string, fields ...zap.Field) {
	l.zap.Error(msg, fields...)
}

func (l *logger) Warn(msg string, fields ...zap.Field) {
	l.zap.Warn(msg, fields...)
}

func (l *logger) Fatal(msg string, fields ...zap.Field) {
	l.zap.Fatal(msg, fields...)
}

// Unstructured logging methods (printf-style)
func (l *logger) Infof(template string, args ...interface{}) {
	l.zap.Sugar().Infof(template, args...)
}

func (l *logger) Debugf(template string, args ...interface{}) {
	l.zap.Sugar().Debugf(template, args...)
}

func (l *logger) Errorf(template string, args ...interface{}) {
	l.zap.Sugar().Errorf(template, args...)
}

func (l *logger) Warnf(template string, args ...interface{}) {
	l.zap.Sugar().Warnf(template, args...)
}

func (l *logger) Fatalf(template string, args ...interface{}) {
	l.zap.Sugar().Fatalf(template, args...)
}

// Utility methods
func (l *logger) Sync() error {
	return l.zap.Sync()
}

func (l *logger) With(fields ...zap.Field) Logger {
	return &logger{zap: l.zap.With(fields...)}
}

// Convenience functions for common structured fields
func String(key, val string) zap.Field       { return zap.String(key, val) }
func Int(key string, val int) zap.Field      { return zap.Int(key, val) }
func Error(err error) zap.Field              { return zap.Error(err) }
func Duration(key string, val interface{}) zap.Field { return zap.Any(key, val) }
func Any(key string, val interface{}) zap.Field      { return zap.Any(key, val) }