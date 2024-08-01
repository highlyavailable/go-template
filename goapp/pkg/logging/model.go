package logging

type LoggerConfig struct {
	Environment      string // "production" OR "development", default is "development"
	WriteStdout      bool   // Write logs to stdout
	EnableStackTrace bool   // Enable stack trace logging
	MaxSize          int    // Max size in megabytes before log is rotated
	MaxBackups       int    // Max number of old log files to keep
	MaxAge           int    // Max number of days to retain old log files
	Compress         bool   // Compress the rotated log files (generates .gz files)
	AppLogPath       string // Path to the app log file
	ErrLogPath       string // Path to the error log file
}
