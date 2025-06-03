package config

import (
	"path/filepath"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config holds all application configuration
type Config struct {
	App           AppConfig           `envconfig:"GO_APP"`
	Database      DatabaseConfig     `envconfig:"POSTGRES"`
	MSSQL         MSSQLConfig        `envconfig:"MSSQL"`
	Logger        LoggerConfig       `envconfig:"LOGGER"`
	Kafka         KafkaConfig        `envconfig:"KAFKA"`
	Observability ObservabilityConfig `envconfig:"OTEL"`
	HTTPClient    HTTPClientConfig    `envconfig:"HTTP_CLIENT"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name        string `envconfig:"APP_NAME" default:"goapp"`
	ProjectRoot string `envconfig:"PROJECT_ROOT" default:"/Users/PeterWBryant/Repos/go-template"`
	EnvPath     string `envconfig:"ENV_PATH"`
	LogDirPath  string `envconfig:"LOG_DIR_PATH"`
	CertDirPath string `envconfig:"CERT_DIR_PATH"`
	Env         string `envconfig:"ENV" default:"development"`
	Port        int    `envconfig:"PORT" default:"8080"`
}

// DatabaseConfig holds PostgreSQL database configuration
type DatabaseConfig struct {
	Host            string `envconfig:"HOST" default:"localhost"`
	Port            int    `envconfig:"PORT" default:"5432"`
	User            string `envconfig:"USER" default:"postgres"`
	Password        string `envconfig:"PASSWORD" default:""`
	DBName          string `envconfig:"NAME" default:"postgres"`
	SSLMode         string `envconfig:"SSLMODE" default:"disable"`
	MaxOpenConns    int    `envconfig:"MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns    int    `envconfig:"MAX_IDLE_CONNS" default:"10"`
	ConnMaxLifetime int    `envconfig:"CONN_MAX_LIFETIME" default:"5"`   // minutes
	ConnMaxIdleTime int    `envconfig:"CONN_MAX_IDLE_TIME" default:"2"`  // minutes
	LogLevel        string `envconfig:"LOG_LEVEL" default:"warn"`        // silent, error, warn, info
}

// MSSQLConfig holds SQL Server database configuration
type MSSQLConfig struct {
	Host            string `envconfig:"HOST" default:"localhost"`
	Port            int    `envconfig:"PORT" default:"1433"`
	User            string `envconfig:"USER" default:"sa"`
	Password        string `envconfig:"PASSWORD" default:""`
	DBName          string `envconfig:"NAME" default:"master"`
	Instance        string `envconfig:"INSTANCE" default:""`
	Encrypt         bool   `envconfig:"ENCRYPT" default:"true"`
	MaxOpenConns    int    `envconfig:"MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns    int    `envconfig:"MAX_IDLE_CONNS" default:"10"`
	ConnMaxLifetime int    `envconfig:"CONN_MAX_LIFETIME" default:"5"`   // minutes
	ConnMaxIdleTime int    `envconfig:"CONN_MAX_IDLE_TIME" default:"2"`  // minutes
	LogLevel        string `envconfig:"LOG_LEVEL" default:"warn"`        // silent, error, warn, info
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Environment      string `envconfig:"ENVIRONMENT" default:"development"`
	WriteStdout      bool   `envconfig:"WRITE_STDOUT" default:"true"`
	EnableStackTrace bool   `envconfig:"ENABLE_STACK_TRACE" default:"false"`
	MaxSize          int    `envconfig:"MAX_SIZE" default:"1"`
	MaxBackups       int    `envconfig:"MAX_BACKUPS" default:"5"`
	MaxAge           int    `envconfig:"MAX_AGE" default:"30"`
	Compress         bool   `envconfig:"COMPRESS" default:"true"`
	AppLogPath       string `envconfig:"APP_LOG_PATH"`
	ErrLogPath       string `envconfig:"ERR_LOG_PATH"`
}

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Brokers        []string `envconfig:"BROKERS" default:"localhost:9092" split_words:"true"`
	ProducerTopic  string   `envconfig:"PRODUCER_TOPIC" default:"events"`
	ConsumerTopic  string   `envconfig:"CONSUMER_TOPIC" default:"events"`
	ConsumerGroup  string   `envconfig:"CONSUMER_GROUP" default:"goapp-group"`
	ConsumerOffset string   `envconfig:"CONSUMER_OFFSET" default:"oldest"`
}

// ObservabilityConfig holds observability configuration
type ObservabilityConfig struct {
	Enabled     bool   `envconfig:"ENABLED" default:"false"`
	ServiceName string `envconfig:"SERVICE_NAME" default:"goapp"`
	Version     string `envconfig:"VERSION" default:"1.0.0"`
	Environment string `envconfig:"ENVIRONMENT" default:"development"`
}

// HTTPClientConfig holds HTTP client configuration
type HTTPClientConfig struct {
	// Timeouts
	Timeout     time.Duration `envconfig:"TIMEOUT" default:"30s"`
	DialTimeout time.Duration `envconfig:"DIAL_TIMEOUT" default:"10s"`
	TLSTimeout  time.Duration `envconfig:"TLS_TIMEOUT" default:"10s"`

	// Connection pooling
	MaxIdleConns        int           `envconfig:"MAX_IDLE_CONNS" default:"100"`
	MaxIdleConnsPerHost int           `envconfig:"MAX_IDLE_CONNS_PER_HOST" default:"10"`
	IdleConnTimeout     time.Duration `envconfig:"IDLE_CONN_TIMEOUT" default:"90s"`

	// TLS
	InsecureSkipVerify bool   `envconfig:"INSECURE_SKIP_VERIFY" default:"false"`
	CertFile           string `envconfig:"CERT_FILE"`
	KeyFile            string `envconfig:"KEY_FILE"`

	// Proxy
	ProxyURL string `envconfig:"PROXY_URL"`

	// Retry
	MaxRetries   int           `envconfig:"MAX_RETRIES" default:"3"`
	RetryWaitMin time.Duration `envconfig:"RETRY_WAIT_MIN" default:"1s"`
	RetryWaitMax time.Duration `envconfig:"RETRY_WAIT_MAX" default:"30s"`

	// Headers
	UserAgent string            `envconfig:"USER_AGENT" default:"goapp/1.0"`
	Headers   map[string]string `envconfig:"HEADERS"`
}

// Load loads configuration from environment variables
func Load() (Config, error) {
	var cfg Config
	
	// Define prefixes and their corresponding config fields
	prefixConfigs := []struct {
		prefix string
		field  interface{}
	}{
		{"GO_APP", &cfg.App},
		{"POSTGRES", &cfg.Database},
		{"MSSQL", &cfg.MSSQL},
		{"LOGGER", &cfg.Logger},
		{"KAFKA", &cfg.Kafka},
		{"OTEL", &cfg.Observability},
		{"HTTP_CLIENT", &cfg.HTTPClient},
	}
	
	// Process each prefix
	for _, pc := range prefixConfigs {
		if err := envconfig.Process(pc.prefix, pc.field); err != nil {
			return cfg, err
		}
	}
	
	// Set default log paths if not provided
	if cfg.Logger.AppLogPath == "" && cfg.App.LogDirPath != "" {
		cfg.Logger.AppLogPath = filepath.Join(cfg.App.LogDirPath, "app.log")
	}
	if cfg.Logger.ErrLogPath == "" && cfg.App.LogDirPath != "" {
		cfg.Logger.ErrLogPath = filepath.Join(cfg.App.LogDirPath, "error.log")
	}
	
	return cfg, nil
}