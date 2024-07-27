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
	Environment      string
	EnableStackTrace bool
	MaxSize          int
	MaxBackups       int
	MaxAge           int
	Compress         bool
	AppLogPath       string
	ErrLogPath       string
}

// InitLogger initializes the logger corresponding to the environment
// value ENV=production OR development. The logs are written to logs/app.log
// and logs/error.log. The logs are rotated based on the configuration
// provided to lumberjack.Logger.
func InitLogger(newLogger LoggerConfig) {
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

	// Ensure each log file directory exists
	// App log directory
	if _, err := os.Stat(filepath.Dir(newLogger.AppLogPath)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(newLogger.AppLogPath), 0755)
		if err != nil {
			panic(err)
		}
	}
	// Error log directory
	if _, err := os.Stat(filepath.Dir(newLogger.ErrLogPath)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(newLogger.ErrLogPath), 0755)
		if err != nil {
			panic(err)
		}
	}

	// Build the zap logger configuration
	var config zap.Config

	// Set the zap logger default configuration based on the environment
	if newLogger.Environment == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	// Customize the logger configuration
	config.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		FunctionKey:    "function",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// Create new zap logger with synced cores for app and error logs
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(config.EncoderConfig),
			zapcore.AddSync(lumberjackLogger),
			zap.DebugLevel, // Log everything to app.log
		),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(config.EncoderConfig),
			zapcore.AddSync(lumberjackErrorLogger),
			zap.ErrorLevel, // Log only errors to error.log
		),
	)

	if newLogger.EnableStackTrace {
		// Create a new logger with the zap logger configuration
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	} else {
		// Create a new logger with the zap logger configuration
		logger = zap.New(core, zap.AddCaller())
	}

	zap.ReplaceGlobals(logger)        // Replace zap's global logger
	logger.Info("Logger initialized") // Announce the logger is initialized
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Writes number entries + 1 to the Error log to test log rotation
func TestRotation(entries int) {
	Info("Dumping " + string(entries) + " entries to the log")
	for i := 0; i < entries; i++ {
		Error("This is an error message")
	}
}
