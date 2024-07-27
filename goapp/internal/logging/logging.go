package logging

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

type LoggerConfig struct {
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	AppLogPath string
	ErrLogPath string
}

// InitLogger initializes the logger corresponding to the environment
// value ENV=production OR development. The logs are written to logs/app.log
// and logs/error.log. The logs are rotated based on the configuration
// provided to lumberjack.Logger.
func InitLogger(newLogger LoggerConfig) {
	// Set the logger configuration based on the environment
	env := os.Getenv("ENV")
	var config zap.Config
	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		env = "development"
		config = zap.NewDevelopmentConfig()
	}

	// Customize the logger configuration
	config.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "Caller",
		MessageKey:     "msg",
		StacktraceKey:  "", // Preference, remove stack trace information
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	// Ensure the logs directory exists
	if _, err := os.Stat(filepath.Dir(newLogger.AppLogPath)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(newLogger.AppLogPath), 0755)
		if err != nil {
			panic(err)
		}
	}
	if _, err := os.Stat(filepath.Dir(newLogger.ErrLogPath)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(newLogger.ErrLogPath), 0755)
		if err != nil {
			panic(err)
		}
	}

	// Use lumberjack for app log rotation
	lumberjackLogger := &lumberjack.Logger{
		Filename:   newLogger.AppLogPath, // Log file path
		MaxSize:    newLogger.MaxSize,    // Max size in megabytes before log is rotated
		MaxBackups: newLogger.MaxBackups, // Max number of old log files to keep
		MaxAge:     newLogger.MaxAge,     // Max number of days to retain old log files
		Compress:   newLogger.Compress,   // Compress the rotated log files
	}

	// Use lumberjack for error log rotation
	lumberjackErrorLogger := &lumberjack.Logger{
		Filename:   newLogger.ErrLogPath, // Error log file path
		MaxSize:    newLogger.MaxSize,    // Max size in megabytes before log is rotated
		MaxBackups: newLogger.MaxBackups, // Max number of old log files to keep
		MaxAge:     newLogger.MaxAge,     // Max number of days to retain old log files
		Compress:   newLogger.Compress,   // Compress the rotated log files
	}

	// Create a core that writes to both app.log and error.log
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(config.EncoderConfig),
			zapcore.AddSync(lumberjackLogger),
			config.Level,
		),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(config.EncoderConfig),
			zapcore.AddSync(lumberjackErrorLogger),
			zapcore.ErrorLevel,
		),
	)

	logger = zap.New(core)

	zap.RedirectStdLog(logger) // Redirects the standard library's log output to the provided logger
	zap.ReplaceGlobals(logger) // Replace zap's global logger
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}
