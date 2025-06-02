package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPartialsHandler_ActivityFeed(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	container := setupTestContainer(t)
	handler := NewPartialsHandler(container)
	
	// Create test router
	router := gin.New()
	router.GET("/partials/activity-feed", handler.ActivityFeed)
	
	// Create test request
	req, err := http.NewRequest(http.MethodGet, "/partials/activity-feed", nil)
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
	
	// Check that response contains activity items
	body := w.Body.String()
	expectedStrings := []string{
		"John Doe created a new post",
		"Jane Smith commented",
		"Alice Johnson joined",
	}
	
	for _, expected := range expectedStrings {
		if !contains(body, expected) {
			t.Errorf("Expected response to contain '%s'", expected)
		}
	}
}

func TestPartialsHandler_Notifications(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	container := setupTestContainer(t)
	handler := NewPartialsHandler(container)
	
	// Create test router
	router := gin.New()
	router.GET("/partials/notifications", handler.Notifications)
	
	// Create test request
	req, err := http.NewRequest(http.MethodGet, "/partials/notifications", nil)
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
	
	// Check that response contains notifications
	body := w.Body.String()
	expectedStrings := []string{
		"New comment on your post",
		"Post published successfully",
	}
	
	for _, expected := range expectedStrings {
		if !contains(body, expected) {
			t.Errorf("Expected response to contain '%s'", expected)
		}
	}
}

func TestPartialsHandler_UserMenu(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	container := setupTestContainer(t)
	handler := NewPartialsHandler(container)
	
	// Create test router
	router := gin.New()
	router.GET("/partials/user-menu", handler.UserMenu)
	
	// Create test request
	req, err := http.NewRequest(http.MethodGet, "/partials/user-menu", nil)
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
	
	// Check that response contains user menu items
	body := w.Body.String()
	expectedStrings := []string{
		"John Doe",
		"Your Profile",
		"Settings",
		"Sign out",
	}
	
	for _, expected := range expectedStrings {
		if !contains(body, expected) {
			t.Errorf("Expected response to contain '%s'", expected)
		}
	}
}

func TestPartialsHandler_MarkNotificationRead(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	container := setupTestContainer(t)
	handler := NewPartialsHandler(container)
	
	// Create test router
	router := gin.New()
	router.POST("/partials/notifications/:id/read", handler.MarkNotificationRead)
	
	// Create test request
	req, err := http.NewRequest(http.MethodPost, "/partials/notifications/123/read", nil)
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
}