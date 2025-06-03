package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Set test environment variables
	os.Setenv("GO_APP_APP_NAME", "test-app")
	os.Setenv("GO_APP_ENV", "test")
	os.Setenv("GO_APP_PORT", "9000")
	os.Setenv("POSTGRES_HOST", "test-host")
	os.Setenv("POSTGRES_PORT", "5433")
	os.Setenv("POSTGRES_USER", "test-user")
	os.Setenv("LOGGER_ENVIRONMENT", "test")
	os.Setenv("LOGGER_WRITE_STDOUT", "false")
	
	defer func() {
		os.Unsetenv("GO_APP_APP_NAME")
		os.Unsetenv("GO_APP_ENV")
		os.Unsetenv("GO_APP_PORT")
		os.Unsetenv("POSTGRES_HOST")
		os.Unsetenv("POSTGRES_PORT")
		os.Unsetenv("POSTGRES_USER")
		os.Unsetenv("LOGGER_ENVIRONMENT")
		os.Unsetenv("LOGGER_WRITE_STDOUT")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test app config
	if cfg.App.Name != "test-app" {
		t.Errorf("Expected app name 'test-app', got '%s'", cfg.App.Name)
	}
	if cfg.App.Env != "test" {
		t.Errorf("Expected env 'test', got '%s'", cfg.App.Env)
	}
	if cfg.App.Port != 9000 {
		t.Errorf("Expected port 9000, got %d", cfg.App.Port)
	}

	// Test database config
	if cfg.Database.Host != "test-host" {
		t.Errorf("Expected database host 'test-host', got '%s'", cfg.Database.Host)
	}
	if cfg.Database.Port != 5433 {
		t.Errorf("Expected database port 5433, got %d", cfg.Database.Port)
	}
	if cfg.Database.User != "test-user" {
		t.Errorf("Expected database user 'test-user', got '%s'", cfg.Database.User)
	}

	// Test logger config
	if cfg.Logger.Environment != "test" {
		t.Errorf("Expected logger environment 'test', got '%s'", cfg.Logger.Environment)
	}
	if cfg.Logger.WriteStdout != false {
		t.Errorf("Expected logger WriteStdout false, got %t", cfg.Logger.WriteStdout)
	}
}

func TestLoadDefaults(t *testing.T) {
	// Clear any existing environment variables that might interfere
	vars := []string{
		"GO_APP_APP_NAME", "GO_APP_ENV", "GO_APP_PORT",
		"POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_USER",
		"LOGGER_ENVIRONMENT", "LOGGER_WRITE_STDOUT",
	}
	
	for _, v := range vars {
		os.Unsetenv(v)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config with defaults: %v", err)
	}

	// Test default values
	if cfg.App.Name != "goapp" {
		t.Errorf("Expected default app name 'goapp', got '%s'", cfg.App.Name)
	}
	if cfg.App.Env != "development" {
		t.Errorf("Expected default env 'development', got '%s'", cfg.App.Env)
	}
	if cfg.App.Port != 8080 {
		t.Errorf("Expected default port 8080, got %d", cfg.App.Port)
	}
	if cfg.Database.Host != "localhost" {
		t.Errorf("Expected default database host 'localhost', got '%s'", cfg.Database.Host)
	}
	if cfg.Database.Port != 5432 {
		t.Errorf("Expected default database port 5432, got %d", cfg.Database.Port)
	}
}

func TestLoadAppError(t *testing.T) {
	os.Setenv("GO_APP_PORT", "invalid_port")
	defer os.Unsetenv("GO_APP_PORT")
	
	_, err := Load()
	if err == nil {
		t.Error("Expected error when app config is invalid")
	}
}

func TestLoadPostgresError(t *testing.T) {
	os.Setenv("POSTGRES_PORT", "invalid_port")
	defer os.Unsetenv("POSTGRES_PORT")
	
	_, err := Load()
	if err == nil {
		t.Error("Expected error when postgres config is invalid")
	}
}

func TestLoadMSSQLError(t *testing.T) {
	os.Setenv("MSSQL_PORT", "invalid_port")
	defer os.Unsetenv("MSSQL_PORT")
	
	_, err := Load()
	if err == nil {
		t.Error("Expected error when MSSQL config is invalid")
	}
}

func TestLoadLoggerError(t *testing.T) {
	os.Setenv("LOGGER_MAX_SIZE", "invalid_size")
	defer os.Unsetenv("LOGGER_MAX_SIZE")
	
	_, err := Load()
	if err == nil {
		t.Error("Expected error when logger config is invalid")
	}
}

func TestLoadKafkaConfig(t *testing.T) {
	// Test successful kafka config loading
	os.Setenv("KAFKA_BROKERS", "localhost:9092,localhost:9093")
	os.Setenv("KAFKA_PRODUCER_TOPIC", "test-producer")
	os.Setenv("KAFKA_CONSUMER_TOPIC", "test-consumer")
	os.Setenv("KAFKA_CONSUMER_GROUP", "test-group")
	os.Setenv("KAFKA_CONSUMER_OFFSET", "newest")
	
	defer func() {
		os.Unsetenv("KAFKA_BROKERS")
		os.Unsetenv("KAFKA_PRODUCER_TOPIC")
		os.Unsetenv("KAFKA_CONSUMER_TOPIC")
		os.Unsetenv("KAFKA_CONSUMER_GROUP")
		os.Unsetenv("KAFKA_CONSUMER_OFFSET")
	}()
	
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Test kafka config was loaded correctly
	if len(cfg.Kafka.Brokers) != 2 {
		t.Errorf("Expected 2 brokers, got %d", len(cfg.Kafka.Brokers))
	}
	if cfg.Kafka.ProducerTopic != "test-producer" {
		t.Errorf("Expected producer topic 'test-producer', got '%s'", cfg.Kafka.ProducerTopic)
	}
	if cfg.Kafka.ConsumerOffset != "newest" {
		t.Errorf("Expected consumer offset 'newest', got '%s'", cfg.Kafka.ConsumerOffset)
	}
}

func TestLoadObservabilityError(t *testing.T) {
	os.Setenv("OTEL_ENABLED", "invalid_bool")
	defer os.Unsetenv("OTEL_ENABLED")
	
	_, err := Load()
	if err == nil {
		t.Error("Expected error when observability config is invalid")
	}
}

func TestLogPathSetting(t *testing.T) {
	os.Setenv("GO_APP_LOG_DIR_PATH", "/test/logs")
	defer os.Unsetenv("GO_APP_LOG_DIR_PATH")
	
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	if cfg.Logger.AppLogPath != "/test/logs/app.log" {
		t.Errorf("Expected app log path '/test/logs/app.log', got '%s'", cfg.Logger.AppLogPath)
	}
	if cfg.Logger.ErrLogPath != "/test/logs/error.log" {
		t.Errorf("Expected error log path '/test/logs/error.log', got '%s'", cfg.Logger.ErrLogPath)
	}
}

func TestLogPathNotSetWhenAlreadyProvided(t *testing.T) {
	os.Setenv("GO_APP_LOG_DIR_PATH", "/test/logs")
	os.Setenv("LOGGER_APP_LOG_PATH", "/custom/app.log")
	os.Setenv("LOGGER_ERR_LOG_PATH", "/custom/error.log")
	
	defer func() {
		os.Unsetenv("GO_APP_LOG_DIR_PATH")
		os.Unsetenv("LOGGER_APP_LOG_PATH")
		os.Unsetenv("LOGGER_ERR_LOG_PATH")
	}()
	
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Should not override existing paths
	if cfg.Logger.AppLogPath != "/custom/app.log" {
		t.Errorf("Expected app log path '/custom/app.log', got '%s'", cfg.Logger.AppLogPath)
	}
	if cfg.Logger.ErrLogPath != "/custom/error.log" {
		t.Errorf("Expected error log path '/custom/error.log', got '%s'", cfg.Logger.ErrLogPath)
	}
}

func TestLoadHTTPClientConfig(t *testing.T) {
	// Test successful HTTP client config loading
	os.Setenv("HTTP_CLIENT_TIMEOUT", "60s")
	os.Setenv("HTTP_CLIENT_MAX_RETRIES", "5")
	os.Setenv("HTTP_CLIENT_USER_AGENT", "test-agent/2.0")
	os.Setenv("HTTP_CLIENT_INSECURE_SKIP_VERIFY", "true")
	os.Setenv("HTTP_CLIENT_PROXY_URL", "http://proxy.example.com:8080")
	
	defer func() {
		os.Unsetenv("HTTP_CLIENT_TIMEOUT")
		os.Unsetenv("HTTP_CLIENT_MAX_RETRIES")
		os.Unsetenv("HTTP_CLIENT_USER_AGENT")
		os.Unsetenv("HTTP_CLIENT_INSECURE_SKIP_VERIFY")
		os.Unsetenv("HTTP_CLIENT_PROXY_URL")
	}()
	
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Test HTTP client config was loaded correctly
	if cfg.HTTPClient.Timeout != 60*time.Second {
		t.Errorf("Expected timeout 60s, got %v", cfg.HTTPClient.Timeout)
	}
	if cfg.HTTPClient.MaxRetries != 5 {
		t.Errorf("Expected max retries 5, got %d", cfg.HTTPClient.MaxRetries)
	}
	if cfg.HTTPClient.UserAgent != "test-agent/2.0" {
		t.Errorf("Expected user agent 'test-agent/2.0', got '%s'", cfg.HTTPClient.UserAgent)
	}
	if cfg.HTTPClient.InsecureSkipVerify != true {
		t.Errorf("Expected insecure skip verify true, got %t", cfg.HTTPClient.InsecureSkipVerify)
	}
	if cfg.HTTPClient.ProxyURL != "http://proxy.example.com:8080" {
		t.Errorf("Expected proxy URL 'http://proxy.example.com:8080', got '%s'", cfg.HTTPClient.ProxyURL)
	}
}

func TestLoadHTTPClientDefaults(t *testing.T) {
	// Clear any existing HTTP client environment variables
	vars := []string{
		"HTTP_CLIENT_TIMEOUT", "HTTP_CLIENT_DIAL_TIMEOUT", "HTTP_CLIENT_TLS_TIMEOUT",
		"HTTP_CLIENT_MAX_IDLE_CONNS", "HTTP_CLIENT_MAX_RETRIES", "HTTP_CLIENT_USER_AGENT",
	}
	
	for _, v := range vars {
		os.Unsetenv(v)
	}
	
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config with defaults: %v", err)
	}
	
	// Test default values
	if cfg.HTTPClient.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout 30s, got %v", cfg.HTTPClient.Timeout)
	}
	if cfg.HTTPClient.DialTimeout != 10*time.Second {
		t.Errorf("Expected default dial timeout 10s, got %v", cfg.HTTPClient.DialTimeout)
	}
	if cfg.HTTPClient.MaxIdleConns != 100 {
		t.Errorf("Expected default max idle conns 100, got %d", cfg.HTTPClient.MaxIdleConns)
	}
	if cfg.HTTPClient.MaxRetries != 3 {
		t.Errorf("Expected default max retries 3, got %d", cfg.HTTPClient.MaxRetries)
	}
	if cfg.HTTPClient.UserAgent != "goapp/1.0" {
		t.Errorf("Expected default user agent 'goapp/1.0', got '%s'", cfg.HTTPClient.UserAgent)
	}
}

func TestLoadHTTPClientError(t *testing.T) {
	os.Setenv("HTTP_CLIENT_MAX_RETRIES", "invalid_number")
	defer os.Unsetenv("HTTP_CLIENT_MAX_RETRIES")
	
	_, err := Load()
	if err == nil {
		t.Error("Expected error when HTTP client config is invalid")
	}
}