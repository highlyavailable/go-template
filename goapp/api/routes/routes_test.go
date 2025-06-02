package routes

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goapp/internal/config"
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
type mockDatabase struct{}

func (m *mockDatabase) DB() *gorm.DB                                                { return nil }
func (m *mockDatabase) Close() error                                               { return nil }
func (m *mockDatabase) Ping(ctx context.Context) error                             { return nil }
func (m *mockDatabase) AutoMigrate(ctx context.Context, dst ...interface{}) error { return nil }
func (m *mockDatabase) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return nil
}
func (m *mockDatabase) WithContext(ctx context.Context) postgres.Database { return m }
func (m *mockDatabase) Health(ctx context.Context) error                  { return nil }

func TestSetupRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockLog := &mockLogger{}
	mockDB := &mockDatabase{}
	
	cfg := config.Config{
		App: config.AppConfig{
			Name: "test-app",
			Env:  "test",
			Port: 8080,
		},
	}
	
	container := &container.Container{
		Config:   cfg,
		Logger:   mockLog,
		Database: mockDB,
	}
	
	router := SetupRouter(container)
	
	if router == nil {
		t.Fatal("Expected router to be non-nil")
	}
	
	// Test that routes are properly registered
	routes := router.Routes()
	
	expectedRoutes := map[string]string{
		"/health":        "GET",
		"/metrics":       "GET",
		"/swagger/*any":  "GET",
	}
	
	foundRoutes := make(map[string]string)
	for _, route := range routes {
		foundRoutes[route.Path] = route.Method
	}
	
	for path, method := range expectedRoutes {
		if foundMethod, exists := foundRoutes[path]; !exists {
			t.Errorf("Expected route %s %s to be registered", method, path)
		} else if foundMethod != method {
			t.Errorf("Expected route %s to have method %s, got %s", path, method, foundMethod)
		}
	}
}

func TestHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockLog := &mockLogger{}
	mockDB := &mockDatabase{}
	
	cfg := config.Config{
		App: config.AppConfig{
			Name: "test-app",
			Env:  "test",
			Port: 8080,
		},
	}
	
	container := &container.Container{
		Config:   cfg,
		Logger:   mockLog,
		Database: mockDB,
	}
	
	router := SetupRouter(container)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestMetricsEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockLog := &mockLogger{}
	mockDB := &mockDatabase{}
	
	cfg := config.Config{
		App: config.AppConfig{
			Name: "test-app",
			Env:  "test",
			Port: 8080,
		},
	}
	
	container := &container.Container{
		Config:   cfg,
		Logger:   mockLog,
		Database: mockDB,
	}
	
	router := SetupRouter(container)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/metrics", nil)
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	// Check that prometheus metrics are in response
	if len(w.Body.String()) == 0 {
		t.Error("Expected metrics response to have content")
	}
}

func TestSwaggerEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockLog := &mockLogger{}
	mockDB := &mockDatabase{}
	
	cfg := config.Config{
		App: config.AppConfig{
			Name: "test-app",
			Env:  "test",
			Port: 8080,
		},
	}
	
	container := &container.Container{
		Config:   cfg,
		Logger:   mockLog,
		Database: mockDB,
	}
	
	router := SetupRouter(container)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/swagger/index.html", nil)
	router.ServeHTTP(w, req)
	
	// Note: This might return 404 if swagger docs aren't generated
	// but the route should still be registered
	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d or %d, got %d", http.StatusOK, http.StatusNotFound, w.Code)
	}
}