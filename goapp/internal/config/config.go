package config

import (
	"path/filepath"

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

// Load loads configuration from environment variables
func Load() (Config, error) {
	var cfg Config
	
	// Load app config
	if err := envconfig.Process("GO_APP", &cfg.App); err != nil {
		return cfg, err
	}
	
	// Load database config
	if err := envconfig.Process("POSTGRES", &cfg.Database); err != nil {
		return cfg, err
	}
	
	// Load MSSQL config
	if err := envconfig.Process("MSSQL", &cfg.MSSQL); err != nil {
		return cfg, err
	}
	
	// Load logger config
	if err := envconfig.Process("LOGGER", &cfg.Logger); err != nil {
		return cfg, err
	}
	
	// Load kafka config  
	if err := envconfig.Process("KAFKA", &cfg.Kafka); err != nil {
		return cfg, err
	}
	
	// Load observability config
	if err := envconfig.Process("OTEL", &cfg.Observability); err != nil {
		return cfg, err
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