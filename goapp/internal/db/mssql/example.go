package mssql

import (
	"context"
	"log"
	"time"

	"goapp/internal/config"
	"gorm.io/gorm"
)

// ExampleNew demonstrates how to use the New function to create a SQL Server client
func ExampleNew() {
	cfg := config.MSSQLConfig{
		Host:     "localhost",
		Port:     1433,
		User:     "sa",
		Password: "YourPassword123!",
		DBName:   "MyDatabase",
		Instance: "", // Optional: named instance
		Encrypt:  true,
	}

	db, err := New(cfg)
	if err != nil {
		log.Fatalf("Could not connect to SQL Server: %v", err)
	}
	defer db.Close()

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	err = db.Ping(ctx)
	if err != nil {
		log.Fatalf("Could not ping SQL Server: %v", err)
	}

	log.Println("Successfully connected to SQL Server")

	// Example query using GORM
	var version string
	err = db.DB().WithContext(ctx).Raw("SELECT @@VERSION as Version").Scan(&version).Error
	if err != nil {
		log.Fatalf("Could not execute query: %v", err)
	}
	log.Printf("SQL Server Version: %s", version)

	// Example stored procedure execution
	err = db.ExecuteProc(ctx, "sp_who2", map[string]interface{}{
		"loginname": "sa",
	})
	if err != nil {
		log.Printf("Could not execute stored procedure: %v", err)
	} else {
		log.Println("Stored procedure executed successfully")
	}

	// Example transaction using GORM
	err = db.Transaction(ctx, func(tx *gorm.DB) error {
		// Do some work in transaction
		var result int
		if err := tx.Raw("SELECT 1").Scan(&result).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Transaction failed: %v", err)
	}

	log.Println("Transaction completed successfully")
}