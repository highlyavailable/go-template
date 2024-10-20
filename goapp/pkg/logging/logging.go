package logging

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var zapLog *zap.Logger

// InitLogger initializes the zapLog zap.Logger corresponding to the environment
// value ENV=production OR development. The LoggerConfig struct is used to configure
// rotation and retention of log files as well as its format and output destination.
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

	// Build the zap Logger configuration
	var config zap.Config

	// Set the zap Logger default configuration based on the environment
	if newLogger.Environment == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	// Customize the Logger configuration
	config.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "Logger",
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

	// Create new zap Logger with synced cores for app and error logs
	appLogCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(config.EncoderConfig),
		zapcore.AddSync(lumberjackLogger),
		zap.DebugLevel, // Log everything to app.log
	)

	errorLogCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(config.EncoderConfig),
		zapcore.AddSync(lumberjackErrorLogger),
		zap.ErrorLevel, // Log only errors to error.log
	)

	var cores []zapcore.Core
	cores = append(cores, appLogCore, errorLogCore)

	// Conditionally add the stdout core
	if newLogger.WriteStdout {
		stdoutCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(config.EncoderConfig),
			zapcore.AddSync(os.Stdout),
			zap.DebugLevel, // Log everything to stdout
		)
		cores = append(cores, stdoutCore)
	}

	core := zapcore.NewTee(cores...)

	if newLogger.EnableStackTrace {
		// Create a new Logger with the zap Logger configuration
		zapLog = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel))
	} else {
		// Create a new Logger with the zap Logger configuration
		zapLog = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	}

	zap.ReplaceGlobals(zapLog) // Replace zap's global Logger
	Info("Logger config", zap.Any("config", newLogger))
}

func Info(message string, fields ...zap.Field) {
	zapLog.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	zapLog.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	zapLog.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	zapLog.Fatal(message, fields...)
}

func Printf(message string, args ...interface{}) {
	zapLog.Info(fmt.Sprintf(message, args...))
}

// Writes number entries + 1 to the Error log to test log rotation
func TestRotation(entries int) {
	zapLog.Info(fmt.Sprintf("Dumping %d entries to the log", entries))
	for i := 0; i < entries; i++ {
		zapLog.Error("This is an error message")
	}
}
