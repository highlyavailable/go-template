package postgres

import "github.com/jmoiron/sqlx"

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
