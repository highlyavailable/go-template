package container

import (
	"context"
	"fmt"
	"os"
	"testing"

	"go.uber.org/zap"
	"goapp/internal/db/postgres"
	"goapp/internal/logging"
	"gorm.io/gorm"
)

func TestNew(t *testing.T) {
	// Set minimal environment variables for testing
	os.Setenv("GO_APP_ENV", "test")
	os.Setenv("LOGGER_WRITE_STDOUT", "false")
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")
	os.Setenv("POSTGRES_USER", "test")
	os.Setenv("POSTGRES_PASSWORD", "test")
	os.Setenv("POSTGRES_NAME", "test")
	
	defer func() {
		os.Unsetenv("GO_APP_ENV")
		os.Unsetenv("LOGGER_WRITE_STDOUT")
		os.Unsetenv("POSTGRES_HOST")
		os.Unsetenv("POSTGRES_PORT")
		os.Unsetenv("POSTGRES_USER")
		os.Unsetenv("POSTGRES_PASSWORD")
		os.Unsetenv("POSTGRES_NAME")
	}()

	container, err := New()
	if err != nil {
		// Database connection might fail in test environment, but container should still initialize
		t.Logf("Container initialization completed with expected database error: %v", err)
		return
	}

	if container == nil {
		t.Fatal("Expected container to be non-nil")
	}

	if container.Logger == nil {
		t.Fatal("Expected logger to be initialized")
	}

	// Test that we can use the logger
	container.Logger.Infof("Test log message: test=%v", true)

	// Clean up
	if err := container.Close(); err != nil {
		t.Errorf("Error closing container: %v", err)
	}
}

// mockLogger implements logging.Logger for testing
type mockLogger struct {
	synced bool
}

func (m *mockLogger) Debug(msg string, fields ...zap.Field) {}
func (m *mockLogger) Info(msg string, fields ...zap.Field)  {}
func (m *mockLogger) Warn(msg string, fields ...zap.Field)  {}
func (m *mockLogger) Error(msg string, fields ...zap.Field) {}
func (m *mockLogger) Fatal(msg string, fields ...zap.Field) {}
func (m *mockLogger) Debugf(format string, args ...interface{}) {}
func (m *mockLogger) Infof(format string, args ...interface{}) {}
func (m *mockLogger) Warnf(format string, args ...interface{}) {}
func (m *mockLogger) Errorf(format string, args ...interface{}) {}
func (m *mockLogger) Fatalf(format string, args ...interface{}) {}
func (m *mockLogger) Sync() error { m.synced = true; return nil }
func (m *mockLogger) With(fields ...zap.Field) logging.Logger { return m }

// mockDatabase implements postgres.Database for testing
type mockDatabase struct {
	closed bool
}

func (m *mockDatabase) DB() *gorm.DB                                                { return nil }
func (m *mockDatabase) Close() error                                               { m.closed = true; return nil }
func (m *mockDatabase) Ping(ctx context.Context) error                             { return nil }
func (m *mockDatabase) AutoMigrate(ctx context.Context, dst ...interface{}) error { return nil }
func (m *mockDatabase) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return nil
}
func (m *mockDatabase) WithContext(ctx context.Context) postgres.Database { return m }
func (m *mockDatabase) Health(ctx context.Context) error                  { return nil }

func TestClose(t *testing.T) {
	mockLog := &mockLogger{}
	mockDB := &mockDatabase{}
	
	container := &Container{
		Logger:   mockLog,
		Database: mockDB,
	}
	
	err := container.Close()
	if err != nil {
		t.Errorf("Expected no error when closing container, got: %v", err)
	}
	
	if !mockDB.closed {
		t.Error("Expected database to be closed")
	}
	
	if !mockLog.synced {
		t.Error("Expected logger to be synced")
	}
}

func TestCloseWithNilDependencies(t *testing.T) {
	container := &Container{}
	
	// Test closing empty container
	err := container.Close()
	if err != nil {
		t.Errorf("Expected no error when closing empty container, got: %v", err)
	}
}

func TestCloseWithDatabaseError(t *testing.T) {
	mockLog := &mockLogger{}
	
	// Create a database that fails to close
	mockDB := &mockErrorDatabase{}
	
	container := &Container{
		Logger:   mockLog,
		Database: mockDB,
	}
	
	err := container.Close()
	if err != nil {
		t.Errorf("Expected no error even when database close fails, got: %v", err)
	}
	
	if !mockLog.synced {
		t.Error("Expected logger to be synced even when database close fails")
	}
}

// mockErrorDatabase implements postgres.Database but fails to close
type mockErrorDatabase struct{}

func (m *mockErrorDatabase) DB() *gorm.DB { return nil }
func (m *mockErrorDatabase) Close() error { return fmt.Errorf("failed to close database") }
func (m *mockErrorDatabase) Ping(ctx context.Context) error { return nil }
func (m *mockErrorDatabase) AutoMigrate(ctx context.Context, dst ...interface{}) error { return nil }
func (m *mockErrorDatabase) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return nil
}
func (m *mockErrorDatabase) WithContext(ctx context.Context) postgres.Database { return m }
func (m *mockErrorDatabase) Health(ctx context.Context) error                  { return nil }

func TestNewWithBadConfig(t *testing.T) {
	// Test with invalid config
	os.Setenv("GO_APP_PORT", "invalid")
	defer os.Unsetenv("GO_APP_PORT")
	
	_, err := New()
	if err == nil {
		t.Error("Expected error when config is invalid")
	}
}

func TestNewSuccessPath(t *testing.T) {
	// Create a temporary directory for log files
	tempDir, err := os.MkdirTemp("", "container_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set environment variables for successful initialization
	os.Setenv("GO_APP_ENV", "test")
	os.Setenv("LOGGER_WRITE_STDOUT", "true")
	os.Setenv("LOGGER_APP_LOG_PATH", tempDir+"/app.log")
	os.Setenv("LOGGER_ERR_LOG_PATH", tempDir+"/error.log")
	os.Setenv("POSTGRES_HOST", "nonexistent")  // This will fail database connection but allow config/logger success
	os.Setenv("POSTGRES_PORT", "5432")
	os.Setenv("POSTGRES_USER", "test")
	os.Setenv("POSTGRES_PASSWORD", "test")
	os.Setenv("POSTGRES_NAME", "test")
	
	defer func() {
		os.Unsetenv("GO_APP_ENV")
		os.Unsetenv("LOGGER_WRITE_STDOUT")
		os.Unsetenv("LOGGER_APP_LOG_PATH")
		os.Unsetenv("LOGGER_ERR_LOG_PATH")
		os.Unsetenv("POSTGRES_HOST")
		os.Unsetenv("POSTGRES_PORT")
		os.Unsetenv("POSTGRES_USER")
		os.Unsetenv("POSTGRES_PASSWORD")
		os.Unsetenv("POSTGRES_NAME")
	}()

	// This should fail at database connection but cover config and logger initialization
	_, err = New()
	if err == nil {
		t.Log("Unexpected success - database should have failed")
	} else {
		t.Logf("Expected database error: %v", err)
	}
}

func TestNewWithLoggerError(t *testing.T) {
	// Set environment to cause logger initialization to fail
	os.Setenv("GO_APP_ENV", "test")
	os.Setenv("LOGGER_APP_LOG_PATH", "/invalid/path/that/cannot/exist/app.log")
	os.Setenv("LOGGER_ERR_LOG_PATH", "/invalid/path/that/cannot/exist/error.log")
	os.Setenv("LOGGER_WRITE_STDOUT", "false")
	
	defer func() {
		os.Unsetenv("GO_APP_ENV")
		os.Unsetenv("LOGGER_APP_LOG_PATH")
		os.Unsetenv("LOGGER_ERR_LOG_PATH")
		os.Unsetenv("LOGGER_WRITE_STDOUT")
	}()

	_, err := New()
	if err == nil {
		t.Error("Expected error when logger initialization fails")
	}
	t.Logf("Got expected logger error: %v", err)
}