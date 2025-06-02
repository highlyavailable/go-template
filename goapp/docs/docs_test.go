package docs

import (
	"testing"

	"github.com/swaggo/swag"
)

func TestSwaggerInfo(t *testing.T) {
	// Test that SwaggerInfo is properly initialized
	if SwaggerInfo == nil {
		t.Fatal("SwaggerInfo should not be nil")
	}

	// Test expected values
	if SwaggerInfo.Version != "1.0" {
		t.Errorf("Expected version '1.0', got '%s'", SwaggerInfo.Version)
	}

	if SwaggerInfo.Host != "localhost:8080" {
		t.Errorf("Expected host 'localhost:8080', got '%s'", SwaggerInfo.Host)
	}

	if SwaggerInfo.BasePath != "/" {
		t.Errorf("Expected base path '/', got '%s'", SwaggerInfo.BasePath)
	}

	if SwaggerInfo.Title != "GoApp REST API" {
		t.Errorf("Expected title 'GoApp REST API', got '%s'", SwaggerInfo.Title)
	}

	if SwaggerInfo.Description != "Production-ready Go REST API with dependency injection" {
		t.Errorf("Expected specific description, got '%s'", SwaggerInfo.Description)
	}

	if SwaggerInfo.InfoInstanceName != "swagger" {
		t.Errorf("Expected instance name 'swagger', got '%s'", SwaggerInfo.InfoInstanceName)
	}

	if SwaggerInfo.SwaggerTemplate != docTemplate {
		t.Error("Expected SwaggerTemplate to match docTemplate")
	}
}

func TestSwaggerRegistration(t *testing.T) {
	// Test that swagger is properly registered
	// The init() function should have registered the swagger info
	
	// Try to get the registered swagger info
	spec := swag.GetSwagger(SwaggerInfo.InstanceName())
	if spec == nil {
		t.Error("Expected swagger spec to be registered, but got nil")
	}
	
	// Verify it's the same instance
	if spec != SwaggerInfo {
		t.Error("Expected registered spec to be the same as SwaggerInfo")
	}
}

func TestDocTemplate(t *testing.T) {
	// Test that docTemplate is not empty
	if docTemplate == "" {
		t.Error("docTemplate should not be empty")
	}

	// Test that it contains expected swagger elements
	expectedElements := []string{
		"swagger",
		"info",
		"paths",
		"/health",
		"schemes",
		"host",
		"basePath",
	}

	for _, element := range expectedElements {
		if !contains(docTemplate, element) {
			t.Errorf("docTemplate should contain '%s'", element)
		}
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}