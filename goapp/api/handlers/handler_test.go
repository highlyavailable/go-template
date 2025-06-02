package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goapp/internal/container"
	"goapp/internal/db/postgres"
	"goapp/internal/logging"
	"gorm.io/gorm"
)

// mockLogger implements logging.Logger for testing
type mockLogger struct {
	logs []string
}

func (m *mockLogger) Debug(msg string, fields ...zap.Field) { m.logs = append(m.logs, "DEBUG: "+msg) }
func (m *mockLogger) Info(msg string, fields ...zap.Field)  { m.logs = append(m.logs, "INFO: "+msg) }
func (m *mockLogger) Warn(msg string, fields ...zap.Field)  { m.logs = append(m.logs, "WARN: "+msg) }
func (m *mockLogger) Error(msg string, fields ...zap.Field) { m.logs = append(m.logs, "ERROR: "+msg) }
func (m *mockLogger) Fatal(msg string, fields ...zap.Field) { m.logs = append(m.logs, "FATAL: "+msg) }

func (m *mockLogger) Debugf(format string, args ...interface{}) {
	m.logs = append(m.logs, "DEBUG: "+fmt.Sprintf(format, args...))
}
func (m *mockLogger) Infof(format string, args ...interface{}) {
	m.logs = append(m.logs, "INFO: "+fmt.Sprintf(format, args...))
}
func (m *mockLogger) Warnf(format string, args ...interface{}) {
	m.logs = append(m.logs, "WARN: "+fmt.Sprintf(format, args...))
}
func (m *mockLogger) Errorf(format string, args ...interface{}) {
	m.logs = append(m.logs, "ERROR: "+fmt.Sprintf(format, args...))
}
func (m *mockLogger) Fatalf(format string, args ...interface{}) {
	m.logs = append(m.logs, "FATAL: "+fmt.Sprintf(format, args...))
}

func (m *mockLogger) Sync() error { return nil }
func (m *mockLogger) With(fields ...zap.Field) logging.Logger {
	return &mockLogger{logs: m.logs}
}

// mockDatabase implements postgres.Database for testing
type mockDatabase struct {
	shouldFailPing   bool
	shouldFailHealth bool
}

func (m *mockDatabase) DB() *gorm.DB                                                { return nil }
func (m *mockDatabase) Close() error                                               { return nil }
func (m *mockDatabase) AutoMigrate(ctx context.Context, dst ...interface{}) error { return nil }
func (m *mockDatabase) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return nil
}
func (m *mockDatabase) WithContext(ctx context.Context) postgres.Database { return m }

func (m *mockDatabase) Ping(ctx context.Context) error {
	if m.shouldFailPing {
		return fmt.Errorf("ping failed")
	}
	return nil
}

func (m *mockDatabase) Health(ctx context.Context) error {
	if m.shouldFailHealth {
		return fmt.Errorf("health check failed")
	}
	return nil
}

func TestNew(t *testing.T) {
	mockLog := &mockLogger{}
	mockDB := &mockDatabase{}
	
	c := &container.Container{
		Logger:   mockLog,
		Database: mockDB,
	}

	handler := New(c)
	if handler == nil {
		t.Fatal("Expected handler to be non-nil")
	}

	if handler.Logger != c.Logger {
		t.Error("Expected handler.Logger to match container.Logger")
	}

	if handler.Database != c.Database {
		t.Error("Expected handler.Database to match container.Database")
	}
}

func TestHealthCheckHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockLog := &mockLogger{}
	mockDB := &mockDatabase{shouldFailPing: false}
	
	handler := &Handler{
		Logger:   mockLog,
		Database: mockDB,
	}
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/health", nil)
	c.Request = req
	
	handler.HealthCheckHandler(c)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if response["status"] != "UP" {
		t.Errorf("Expected status UP, got %s", response["status"])
	}
	
	// Check that logger was called
	if len(mockLog.logs) == 0 {
		t.Error("Expected logger to be called")
	}
}

func TestHealthCheckHandler_DatabaseFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockLog := &mockLogger{}
	mockDB := &mockDatabase{shouldFailPing: true}
	
	handler := &Handler{
		Logger:   mockLog,
		Database: mockDB,
	}
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/health", nil)
	c.Request = req
	
	handler.HealthCheckHandler(c)
	
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
	}
	
	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if response["status"] != "DOWN" {
		t.Errorf("Expected status DOWN, got %s", response["status"])
	}
	
	if response["error"] != "Database connection failed" {
		t.Errorf("Expected error message 'Database connection failed', got %s", response["error"])
	}
}