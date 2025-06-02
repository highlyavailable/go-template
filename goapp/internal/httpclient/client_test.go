package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"goapp/internal/config"
	"goapp/internal/logging"
)

func TestNew(t *testing.T) {
	cfg := Config{
		Timeout:         30 * time.Second,
		MaxRetries:      3,
		UserAgent:       "test-agent",
		MaxIdleConns:    10,
	}
	
	// Create a mock logger for testing
	loggerCfg := config.LoggerConfig{
		Environment: "test",
		WriteStdout: false,
	}
	logger, err := logging.New(loggerCfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	client, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}
	
	if client.config.UserAgent != "test-agent" {
		t.Errorf("Expected user agent 'test-agent', got '%s'", client.config.UserAgent)
	}
}

func TestClient_Get(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()
	
	// Create client
	cfg := Config{
		Timeout:    5 * time.Second,
		MaxRetries: 1,
	}
	
	loggerCfg := config.LoggerConfig{
		Environment: "test",
		WriteStdout: false,
	}
	logger, _ := logging.New(loggerCfg)
	
	client, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	// Test GET request
	ctx := context.Background()
	resp, err := client.Get(ctx, server.URL)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestClient_Retry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()
	
	cfg := Config{
		Timeout:      5 * time.Second,
		MaxRetries:   3,
		RetryWaitMin: 10 * time.Millisecond,
		RetryWaitMax: 100 * time.Millisecond,
	}
	
	loggerCfg := config.LoggerConfig{
		Environment: "test",
		WriteStdout: false,
	}
	logger, _ := logging.New(loggerCfg)
	
	client, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	ctx := context.Background()
	resp, err := client.Get(ctx, server.URL)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestClient_HealthCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	cfg := Config{
		Timeout: 5 * time.Second,
	}
	
	loggerCfg := config.LoggerConfig{
		Environment: "test",
		WriteStdout: false,
	}
	logger, _ := logging.New(loggerCfg)
	
	client, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	ctx := context.Background()
	err = client.HealthCheck(ctx, server.URL)
	if err != nil {
		t.Errorf("Health check failed: %v", err)
	}
}

func TestClient_HealthCheckFail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	
	cfg := Config{
		Timeout: 5 * time.Second,
	}
	
	loggerCfg := config.LoggerConfig{
		Environment: "test",
		WriteStdout: false,
	}
	logger, _ := logging.New(loggerCfg)
	
	client, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	ctx := context.Background()
	err = client.HealthCheck(ctx, server.URL)
	if err == nil {
		t.Error("Expected health check to fail, but it succeeded")
	}
}

func TestDecodeJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name":"test","value":42}`))
	}))
	defer server.Close()
	
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	
	var result struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	
	err = DecodeJSON(resp, &result)
	if err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}
	
	if result.Name != "test" {
		t.Errorf("Expected name 'test', got '%s'", result.Name)
	}
	
	if result.Value != 42 {
		t.Errorf("Expected value 42, got %d", result.Value)
	}
}

func TestClient_Post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status":"created"}`))
	}))
	defer server.Close()
	
	cfg := Config{
		Timeout:    5 * time.Second,
		MaxRetries: 1,
	}
	
	loggerCfg := config.LoggerConfig{
		Environment: "test",
		WriteStdout: false,
	}
	logger, _ := logging.New(loggerCfg)
	
	client, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	ctx := context.Background()
	data := map[string]string{"key": "value"}
	resp, err := client.Post(ctx, server.URL, "application/json", data)
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}
}

func TestClient_Put(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	cfg := Config{
		Timeout:    5 * time.Second,
		MaxRetries: 1,
	}
	
	loggerCfg := config.LoggerConfig{
		Environment: "test",
		WriteStdout: false,
	}
	logger, _ := logging.New(loggerCfg)
	
	client, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	ctx := context.Background()
	data := map[string]string{"key": "value"}
	resp, err := client.Put(ctx, server.URL, "application/json", data)
	if err != nil {
		t.Fatalf("PUT request failed: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	
	cfg := Config{
		Timeout:    5 * time.Second,
		MaxRetries: 1,
	}
	
	loggerCfg := config.LoggerConfig{
		Environment: "test",
		WriteStdout: false,
	}
	logger, _ := logging.New(loggerCfg)
	
	client, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	ctx := context.Background()
	resp, err := client.Delete(ctx, server.URL)
	if err != nil {
		t.Fatalf("DELETE request failed: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}

func TestClient_Close(t *testing.T) {
	cfg := Config{
		Timeout:    5 * time.Second,
		MaxRetries: 1,
	}
	
	loggerCfg := config.LoggerConfig{
		Environment: "test",
		WriteStdout: false,
	}
	logger, _ := logging.New(loggerCfg)
	
	client, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	// Should not panic
	client.Close()
}

func TestClient_WithCustomHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom-Header") != "test-value" {
			t.Errorf("Expected custom header 'test-value', got '%s'", r.Header.Get("X-Custom-Header"))
		}
		if r.Header.Get("User-Agent") != "custom-agent" {
			t.Errorf("Expected user agent 'custom-agent', got '%s'", r.Header.Get("User-Agent"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	cfg := Config{
		Timeout:   5 * time.Second,
		UserAgent: "custom-agent",
		Headers: map[string]string{
			"X-Custom-Header": "test-value",
		},
	}
	
	loggerCfg := config.LoggerConfig{
		Environment: "test",
		WriteStdout: false,
	}
	logger, _ := logging.New(loggerCfg)
	
	client, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	ctx := context.Background()
	resp, err := client.Get(ctx, server.URL)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	defer resp.Body.Close()
}

func TestClient_MaxRetries(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusInternalServerError) // Always fail
	}))
	defer server.Close()
	
	cfg := Config{
		Timeout:      1 * time.Second,
		MaxRetries:   2,
		RetryWaitMin: 10 * time.Millisecond,
		RetryWaitMax: 100 * time.Millisecond,
	}
	
	loggerCfg := config.LoggerConfig{
		Environment: "test",
		WriteStdout: false,
	}
	logger, _ := logging.New(loggerCfg)
	
	client, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	ctx := context.Background()
	resp, err := client.Get(ctx, server.URL)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	defer resp.Body.Close()
	
	// Should have attempted 3 times total (initial + 2 retries)
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestNew_WithInvalidProxy(t *testing.T) {
	cfg := Config{
		Timeout:   30 * time.Second,
		ProxyURL:  "://invalid-url", // Invalid URL without scheme
	}
	
	loggerCfg := config.LoggerConfig{
		Environment: "test",
		WriteStdout: false,
	}
	logger, _ := logging.New(loggerCfg)
	
	_, err := New(cfg, logger)
	if err == nil {
		t.Error("Expected error with invalid proxy URL, but got none")
	}
}

func TestNew_WithValidProxy(t *testing.T) {
	cfg := Config{
		Timeout:   30 * time.Second,
		ProxyURL:  "http://proxy.example.com:8080",
	}
	
	loggerCfg := config.LoggerConfig{
		Environment: "test",
		WriteStdout: false,
	}
	logger, _ := logging.New(loggerCfg)
	
	client, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create client with valid proxy: %v", err)
	}
	
	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}
}

func TestNew_WithInvalidCertificates(t *testing.T) {
	cfg := Config{
		Timeout:  30 * time.Second,
		CertFile: "/nonexistent/cert.pem",
		KeyFile:  "/nonexistent/key.pem",
	}
	
	loggerCfg := config.LoggerConfig{
		Environment: "test",
		WriteStdout: false,
	}
	logger, _ := logging.New(loggerCfg)
	
	_, err := New(cfg, logger)
	if err == nil {
		t.Error("Expected error with invalid certificates, but got none")
	}
}

func TestClient_WithRedirects(t *testing.T) {
	redirectCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if redirectCount < 2 {
			redirectCount++
			// Use proper redirect URL format
			redirectURL := "/redirect" + string(rune('0'+redirectCount))
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("final response"))
	}))
	defer server.Close()

	cfg := Config{
		Timeout:    5 * time.Second,
		MaxRetries: 1,
	}

	loggerCfg := config.LoggerConfig{
		Environment: "test",
		WriteStdout: false,
	}
	logger, _ := logging.New(loggerCfg)

	client, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	resp, err := client.Get(ctx, server.URL)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if redirectCount != 2 {
		t.Errorf("Expected 2 redirects, got %d", redirectCount)
	}
}

func TestClient_TooManyRedirects(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Always redirect to create infinite loop
		http.Redirect(w, r, r.URL.Path+"?redirect=loop", http.StatusFound)
	}))
	defer server.Close()

	cfg := Config{
		Timeout:    5 * time.Second,
		MaxRetries: 0, // No retries for this test
	}

	loggerCfg := config.LoggerConfig{
		Environment: "test",
		WriteStdout: false,
	}
	logger, _ := logging.New(loggerCfg)

	client, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	resp, err := client.Get(ctx, server.URL)
	if err == nil {
		defer resp.Body.Close()
		t.Error("Expected error due to too many redirects, but got none")
	}
}

func TestDecodeJSON_ErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"bad request"}`))
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	var result map[string]interface{}
	err = DecodeJSON(resp, &result)
	if err == nil {
		t.Error("Expected error due to 400 status, but got none")
	}
}

func TestReadBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response body"))
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	body, err := ReadBody(resp)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}

	expected := "test response body"
	if string(body) != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, string(body))
	}
}

func TestReadBodyAsString(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test string response"))
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	body, err := ReadBodyAsString(resp)
	if err != nil {
		t.Fatalf("Failed to read body as string: %v", err)
	}

	expected := "test string response"
	if body != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, body)
	}
}

func TestCreateRequestWithBody_StringBody(t *testing.T) {
	ctx := context.Background()
	body := "test string body"
	
	req, err := createRequestWithBody(ctx, "POST", "http://example.com", "text/plain", body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	if req.Header.Get("Content-Type") != "text/plain" {
		t.Errorf("Expected Content-Type 'text/plain', got '%s'", req.Header.Get("Content-Type"))
	}
}

func TestCreateRequestWithBody_BytesBody(t *testing.T) {
	ctx := context.Background()
	body := []byte("test bytes body")
	
	req, err := createRequestWithBody(ctx, "POST", "http://example.com", "application/octet-stream", body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	if req.Header.Get("Content-Type") != "application/octet-stream" {
		t.Errorf("Expected Content-Type 'application/octet-stream', got '%s'", req.Header.Get("Content-Type"))
	}
}

func TestCreateRequestWithBody_IOReaderBody(t *testing.T) {
	ctx := context.Background()
	body := strings.NewReader("test reader body")
	
	req, err := createRequestWithBody(ctx, "POST", "http://example.com", "text/plain", body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	if req.Header.Get("Content-Type") != "text/plain" {
		t.Errorf("Expected Content-Type 'text/plain', got '%s'", req.Header.Get("Content-Type"))
	}
}

func TestCreateRequestWithBody_JSONBodyWithoutContentType(t *testing.T) {
	ctx := context.Background()
	body := map[string]string{"key": "value"}
	
	req, err := createRequestWithBody(ctx, "POST", "http://example.com", "", body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	
	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", req.Header.Get("Content-Type"))
	}
}

func TestCreateRequestWithBody_InvalidJSON(t *testing.T) {
	ctx := context.Background()
	// Use a channel which cannot be marshaled to JSON
	body := make(chan int)
	
	_, err := createRequestWithBody(ctx, "POST", "http://example.com", "", body)
	if err == nil {
		t.Error("Expected error with invalid JSON body, but got none")
	}
}

func TestCreateProxyFunc(t *testing.T) {
	proxyURL := "http://proxy.example.com:8080"
	
	proxyFunc, err := createProxyFunc(proxyURL)
	if err != nil {
		t.Fatalf("Failed to create proxy function: %v", err)
	}
	
	if proxyFunc == nil {
		t.Fatal("Expected proxy function to be non-nil")
	}
	
	// Test the proxy function with a dummy request
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	url, err := proxyFunc(req)
	if err != nil {
		t.Fatalf("Proxy function failed: %v", err)
	}
	
	if url.String() != proxyURL {
		t.Errorf("Expected proxy URL '%s', got '%s'", proxyURL, url.String())
	}
}

func TestCreateProxyFunc_InvalidURL(t *testing.T) {
	proxyURL := "://invalid-url"
	
	_, err := createProxyFunc(proxyURL)
	if err == nil {
		t.Error("Expected error with invalid proxy URL, but got none")
	}
}