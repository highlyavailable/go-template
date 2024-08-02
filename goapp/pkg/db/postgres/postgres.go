package postgres

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// DB is a global variable to hold the database connection pool
var DB *sqlx.DB

// Config holds the configuration for the database connection
type Config struct {
	Host     string `envconfig:"HOST" default:"localhost"`
	Port     int    `envconfig:"PORT" default:"5432"`
	User     string `envconfig:"USER" default:"postgres"`
	Password string `envconfig:"PASSWORD" default:""`
	DBName   string `envconfig:"NAME" default:"postgres"`
	SSLMode  string `envconfig:"SSLMODE" default:"disable"`
}

// NewConfig creates a new database configuration from environment variables
func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("POSTGRES", &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to process envconfig: %w", err)
	}
	return &cfg, nil
}

// Connect initializes the database connection
func Connect(cfg *Config) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	DB = db
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// Migrate runs database migrations
func Migrate(migrationsDir string) error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	goose.SetDialect("postgres")
	err := goose.Up(DB.DB, migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// Ping checks the database connection
func Ping() error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	return DB.Ping()
}
