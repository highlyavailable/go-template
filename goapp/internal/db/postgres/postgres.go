package postgres

import (
	"context"
	"fmt"
	"time"

	"goapp/internal/config"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database interface defines database operations for PostgreSQL with GORM
type Database interface {
	DB() *gorm.DB
	Close() error
	Ping(ctx context.Context) error
	AutoMigrate(ctx context.Context, dst ...interface{}) error
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	WithContext(ctx context.Context) Database
	Health(ctx context.Context) error
}

// postgres implements the Database interface using GORM
type postgres struct {
	db *gorm.DB
}

// New creates a new PostgreSQL database instance with GORM
func New(cfg config.DatabaseConfig) (Database, error) {
	p := &postgres{}
	if err := p.connect(cfg); err != nil {
		return nil, err
	}
	return p, nil
}

// connect initializes the PostgreSQL database connection
func (p *postgres) connect(cfg config.DatabaseConfig) error {
	// Build DSN for PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	// Configure GORM logger
	gormLogger := logger.Default
	if cfg.LogLevel != "" {
		switch cfg.LogLevel {
		case "silent":
			gormLogger = logger.Default.LogMode(logger.Silent)
		case "error":
			gormLogger = logger.Default.LogMode(logger.Error)
		case "warn":
			gormLogger = logger.Default.LogMode(logger.Warn)
		case "info":
			gormLogger = logger.Default.LogMode(logger.Info)
		}
	}

	// Open database connection
	db, err := gorm.Open(postgresDriver.Open(dsn), &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: false,
		SkipDefaultTransaction:                   false,
		PrepareStmt:                             true,
		QueryFields:                             true,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Get underlying SQL DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Minute)

	p.db = db
	return nil
}

// DB returns the underlying GORM DB instance
func (p *postgres) DB() *gorm.DB {
	return p.db
}

// Close closes the database connection
func (p *postgres) Close() error {
	if p.db == nil {
		return nil
	}
	
	sqlDB, err := p.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	
	return sqlDB.Close()
}

// Ping checks the database connection
func (p *postgres) Ping(ctx context.Context) error {
	if p.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	
	sqlDB, err := p.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	
	return sqlDB.PingContext(ctx)
}

// AutoMigrate runs database migrations for the given models
func (p *postgres) AutoMigrate(ctx context.Context, dst ...interface{}) error {
	if p.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	
	db := p.db.WithContext(ctx)
	return db.AutoMigrate(dst...)
}

// Transaction executes a function within a database transaction
func (p *postgres) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	if p.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	
	return p.db.WithContext(ctx).Transaction(fn)
}

// WithContext returns a new Database instance with the given context
func (p *postgres) WithContext(ctx context.Context) Database {
	if p.db == nil {
		return p
	}
	
	return &postgres{
		db: p.db.WithContext(ctx),
	}
}

// Health performs a comprehensive health check
func (p *postgres) Health(ctx context.Context) error {
	if err := p.Ping(ctx); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	
	// Test a simple query
	var result int
	if err := p.db.WithContext(ctx).Raw("SELECT 1").Scan(&result).Error; err != nil {
		return fmt.Errorf("health query failed: %w", err)
	}
	
	if result != 1 {
		return fmt.Errorf("health query returned unexpected result: %d", result)
	}
	
	return nil
}