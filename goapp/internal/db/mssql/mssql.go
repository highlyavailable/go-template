package mssql

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"goapp/internal/config"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database interface defines database operations for SQL Server with GORM
type Database interface {
	DB() *gorm.DB
	Close() error
	Ping(ctx context.Context) error
	AutoMigrate(ctx context.Context, dst ...interface{}) error
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	WithContext(ctx context.Context) Database
	Health(ctx context.Context) error
	ExecuteProc(ctx context.Context, procName string, params map[string]interface{}) error
}

// mssql implements the Database interface using GORM
type mssql struct {
	db *gorm.DB
}

// New creates a new SQL Server database instance with GORM
func New(cfg config.MSSQLConfig) (Database, error) {
	m := &mssql{}
	if err := m.connect(cfg); err != nil {
		return nil, err
	}
	return m, nil
}

// connect initializes the SQL Server database connection
func (m *mssql) connect(cfg config.MSSQLConfig) error {
	// Build connection string for SQL Server
	connectionString := buildConnectionString(cfg)

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
	db, err := gorm.Open(sqlserver.Open(connectionString), &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: false,
		SkipDefaultTransaction:                   false,
		PrepareStmt:                             true,
		QueryFields:                             true,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to SQL Server: %w", err)
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

	m.db = db
	return nil
}

// buildConnectionString constructs a SQL Server connection string
func buildConnectionString(cfg config.MSSQLConfig) string {
	// Build the base URL for GORM SQL Server driver
	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}

	// Add query parameters
	query := url.Values{}
	
	if cfg.DBName != "" {
		query.Add("database", cfg.DBName)
	}
	
	if cfg.Instance != "" {
		query.Add("instance", cfg.Instance)
	}
	
	// Encryption settings
	if cfg.Encrypt {
		query.Add("encrypt", "true")
		query.Add("TrustServerCertificate", "false")
	} else {
		query.Add("encrypt", "false")
	}
	
	// Connection timeout
	query.Add("connection timeout", "30")
	
	// Enable Multiple Active Result Sets (MARS)
	query.Add("MultipleActiveResultSets", "true")
	
	u.RawQuery = query.Encode()
	
	return u.String()
}

// DB returns the underlying GORM DB instance
func (m *mssql) DB() *gorm.DB {
	return m.db
}

// Close closes the database connection
func (m *mssql) Close() error {
	if m.db == nil {
		return nil
	}
	
	sqlDB, err := m.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	
	return sqlDB.Close()
}

// Ping checks the database connection
func (m *mssql) Ping(ctx context.Context) error {
	if m.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	
	sqlDB, err := m.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	
	return sqlDB.PingContext(ctx)
}

// AutoMigrate runs database migrations for the given models
func (m *mssql) AutoMigrate(ctx context.Context, dst ...interface{}) error {
	if m.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	
	db := m.db.WithContext(ctx)
	return db.AutoMigrate(dst...)
}

// Transaction executes a function within a database transaction
func (m *mssql) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	if m.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	
	return m.db.WithContext(ctx).Transaction(fn)
}

// WithContext returns a new Database instance with the given context
func (m *mssql) WithContext(ctx context.Context) Database {
	if m.db == nil {
		return m
	}
	
	return &mssql{
		db: m.db.WithContext(ctx),
	}
}

// Health performs a comprehensive health check
func (m *mssql) Health(ctx context.Context) error {
	if err := m.Ping(ctx); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	
	// Test a simple query
	var result int
	if err := m.db.WithContext(ctx).Raw("SELECT 1").Scan(&result).Error; err != nil {
		return fmt.Errorf("health query failed: %w", err)
	}
	
	if result != 1 {
		return fmt.Errorf("health query returned unexpected result: %d", result)
	}
	
	return nil
}

// ExecuteProc executes a stored procedure with parameters
func (m *mssql) ExecuteProc(ctx context.Context, procName string, params map[string]interface{}) error {
	if m.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	
	// Build EXEC statement
	query := fmt.Sprintf("EXEC %s", procName)
	args := make([]interface{}, 0, len(params))
	
	if len(params) > 0 {
		paramNames := make([]string, 0, len(params))
		for name, value := range params {
			paramNames = append(paramNames, fmt.Sprintf("@%s = ?", name))
			args = append(args, value)
		}
		query += " " + paramNames[0]
		for i := 1; i < len(paramNames); i++ {
			query += ", " + paramNames[i]
		}
	}
	
	return m.db.WithContext(ctx).Exec(query, args...).Error
}