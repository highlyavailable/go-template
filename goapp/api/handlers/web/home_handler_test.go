package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"goapp/internal/container"
	"goapp/internal/config"
	"goapp/internal/logging"
)

func setupTestContainer(t *testing.T) *container.Container {
	// Create test config
	cfg := config.Config{
		Logger: config.LoggerConfig{
			Environment: "test",
			WriteStdout: true,
		},
	}
	
	// Create test logger
	logger, err := logging.New(cfg.Logger)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	// Return container with minimal dependencies
	return &container.Container{
		Config: cfg,
		Logger: logger,
	}
}

func TestHomeHandler_Index(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	container := setupTestContainer(t)
	handler := NewHomeHandler(container)
	
	// Create test router
	router := gin.New()
	router.GET("/", handler.Index)
	
	// Create test request
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assertions
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	// Check content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html" {
		t.Errorf("Expected content type 'text/html', got '%s'", contentType)
	}
	
	// Check that response body contains expected content
	body := w.Body.String()
	expectedStrings := []string{
		"Dashboard",
		"GoApp",
		"Welcome to GoApp",
	}
	
	for _, expected := range expectedStrings {
		if !contains(body, expected) {
			t.Errorf("Expected response to contain '%s'", expected)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		(s == substr || len(s) > len(substr) && 
		(s[:len(substr)] == substr || 
		contains(s[1:], substr)))
}