package logging

import (
	"os"
	"path/filepath"
	"testing"

	"goapp/internal/config"
)

func TestNew(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "logging_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := config.LoggerConfig{
		Environment:      "test",
		WriteStdout:      false,
		EnableStackTrace: false,
		MaxSize:          1,
		MaxBackups:       5,
		MaxAge:           30,
		Compress:         true,
		AppLogPath:       filepath.Join(tempDir, "app.log"),
		ErrLogPath:       filepath.Join(tempDir, "error.log"),
	}

	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	if logger == nil {
		t.Fatal("Expected logger to be non-nil")
	}

	// Test structured logging
	logger.Info("Test info message", String("key", "value"), Int("count", 42))
	logger.Debug("Test debug message", Any("data", map[string]string{"test": "value"}))
	logger.Error("Test error message", Error(err))
	logger.Warn("Test warn message")

	// Test unstructured logging
	logger.Infof("Test printf message %s %d", "formatted", 123)
	logger.Debugf("Debug: %v", true)
	logger.Errorf("Error: %v", "something went wrong")
	logger.Warnf("Warning: %s", "test warning")

	// Test With method
	contextLogger := logger.With(String("component", "test"))
	contextLogger.Info("Context logger test")

	// Test Sync
	err = logger.Sync()
	if err != nil {
		t.Errorf("Error syncing logger: %v", err)
	}
}

func TestConvenienceFunctions(t *testing.T) {
	// Test convenience field functions
	stringField := String("test", "value")
	if stringField.Key != "test" {
		t.Errorf("Expected key 'test', got '%s'", stringField.Key)
	}

	intField := Int("count", 42)
	if intField.Key != "count" {
		t.Errorf("Expected key 'count', got '%s'", intField.Key)
	}

	testErr := &testError{msg: "test error"}
	errorField := Error(testErr)
	if errorField.Key != "error" {
		t.Errorf("Expected key 'error', got '%s'", errorField.Key)
	}

	anyField := Any("data", map[string]int{"test": 1})
	if anyField.Key != "data" {
		t.Errorf("Expected key 'data', got '%s'", anyField.Key)
	}
}

func TestLoggerInterface(t *testing.T) {
	// Test that our logger implements the Logger interface
	var _ Logger = (*logger)(nil)
}

func TestWithMethod(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "logging_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := config.LoggerConfig{
		Environment: "test",
		WriteStdout: false,
		AppLogPath:  filepath.Join(tempDir, "app.log"),
		ErrLogPath:  filepath.Join(tempDir, "error.log"),
	}

	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test With returns a new logger instance
	contextLogger := logger.With(String("service", "test-service"))
	if contextLogger == logger {
		t.Error("With should return a new logger instance")
	}

	// Test that both loggers work
	logger.Info("Original logger")
	contextLogger.Info("Context logger with service field")
}

func TestDevelopmentLogger(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "logging_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := config.LoggerConfig{
		Environment:      "development",
		WriteStdout:      true,
		EnableStackTrace: true,
		AppLogPath:       filepath.Join(tempDir, "app.log"),
		ErrLogPath:       filepath.Join(tempDir, "error.log"),
	}

	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create development logger: %v", err)
	}

	if logger == nil {
		t.Fatal("Expected logger to be non-nil")
	}

	// Test that development logger works
	logger.Info("Development mode info")
	logger.Debug("Debug message should be visible in development")
}

func TestProductionLogger(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "logging_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := config.LoggerConfig{
		Environment:      "production",
		WriteStdout:      false,
		EnableStackTrace: false,
		AppLogPath:       filepath.Join(tempDir, "app.log"),
		ErrLogPath:       filepath.Join(tempDir, "error.log"),
	}

	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create production logger: %v", err)
	}

	if logger == nil {
		t.Fatal("Expected logger to be non-nil")
	}

	// Test that production logger works
	logger.Info("Production mode info")
	logger.Error("Error in production")
}

func TestLoggerWithoutFilePaths(t *testing.T) {
	cfg := config.LoggerConfig{
		Environment: "test",
		WriteStdout: true,
		// No file paths set - should only log to stdout
	}

	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger without file paths: %v", err)
	}

	if logger == nil {
		t.Fatal("Expected logger to be non-nil")
	}

	// Test that logger works without file outputs
	logger.Info("Stdout only message")
	logger.Error("Error to stdout only")
}

func TestInvalidLogPath(t *testing.T) {
	cfg := config.LoggerConfig{
		Environment: "test",
		WriteStdout: false,
		AppLogPath:  "/invalid/path/that/should/not/exist/app.log",
		ErrLogPath:  "/invalid/path/that/should/not/exist/error.log",
	}

	_, err := New(cfg)
	if err == nil {
		t.Error("Expected error when using invalid log path")
	}
}

func TestLoggerLevels(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "logging_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := config.LoggerConfig{
		Environment: "test",
		WriteStdout: false,
		AppLogPath:  filepath.Join(tempDir, "app.log"),
		ErrLogPath:  filepath.Join(tempDir, "error.log"),
	}

	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test all log levels
	logger.Debug("Debug level message", String("level", "debug"))
	logger.Info("Info level message", String("level", "info"))
	logger.Warn("Warn level message", String("level", "warn"))
	logger.Error("Error level message", String("level", "error"))

	// Test formatted versions
	logger.Debugf("Debug: %s", "formatted")
	logger.Infof("Info: %s", "formatted")
	logger.Warnf("Warn: %s", "formatted")
	logger.Errorf("Error: %s", "formatted")
}

func TestMultipleWith(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "logging_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := config.LoggerConfig{
		Environment: "test",
		WriteStdout: false,
		AppLogPath:  filepath.Join(tempDir, "app.log"),
		ErrLogPath:  filepath.Join(tempDir, "error.log"),
	}

	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test chaining With calls
	contextLogger := logger.With(String("component", "test")).
		With(String("version", "1.0")).
		With(Int("build", 123))

	contextLogger.Info("Message with multiple context fields")
}

// Helper type for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}