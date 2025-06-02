package handlers

import (
	"goapp/internal/container"
	"goapp/internal/db/postgres"
	"goapp/internal/logging"
)

// Handler contains all dependencies for HTTP handlers
type Handler struct {
	Logger   logging.Logger
	Database postgres.Database
}

// New creates a new handler with injected dependencies
func New(container *container.Container) *Handler {
	return &Handler{
		Logger:   container.Logger,
		Database: container.Database,
	}
}