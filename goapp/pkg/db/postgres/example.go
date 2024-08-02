package postgres

import "log"

func ExampleNewConfig() {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatalf("Failed to create database configuration: %v", err)
	}
	err = Connect(cfg)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer Close()
	err = Migrate("migrations")
	if err != nil {
		log.Fatalf("Could not run migrations: %v", err)
	}
	err = Ping()
	if err != nil {
		log.Fatalf("Could not ping the database: %v", err)
	}
}
