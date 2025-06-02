package container

import (
	"goapp/internal/config"
	"goapp/internal/db/postgres"
	"goapp/internal/logging"
	"go.uber.org/zap"
)

// Container holds all application dependencies
type Container struct {
	Config   config.Config
	Logger   logging.Logger
	Database postgres.Database
}

// New creates a new dependency injection container
func New() (*Container, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Initialize logger
	logger, err := logging.New(cfg.Logger)
	if err != nil {
		return nil, err
	}

	// Initialize database
	database, err := postgres.New(cfg.Database)
	if err != nil {
		logger.Warn("Failed to connect to PostgreSQL, continuing without database", zap.Error(err))
		// In a real app, you might want to use SQLite as fallback or fail here
		// For now, we'll continue without database for demo purposes
		database = nil
	}

	return &Container{
		Config:   cfg,
		Logger:   logger,
		Database: database,
	}, nil
}

// Close gracefully shuts down all dependencies
func (c *Container) Close() error {
	if c.Database != nil {
		if err := c.Database.Close(); err != nil {
			c.Logger.Errorf("Failed to close database: %v", err)
		}
	}
	
	if c.Logger != nil {
		c.Logger.Sync()
	}
	
	return nil
}