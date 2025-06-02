package postgres

import (
	"context"
	"fmt"
	"log"

	"goapp/internal/config"
	"goapp/internal/db/migrations"
	"goapp/internal/models"
	"gorm.io/gorm"
)

func ExampleNew() {
	ctx := context.Background()
	
	cfg := config.DatabaseConfig{
		Host:            "localhost",
		Port:            5432,
		User:            "postgres",
		Password:        "password",
		DBName:          "myapp",
		SSLMode:         "disable",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5,
		ConnMaxIdleTime: 1,
		LogLevel:        "info",
	}
	
	db, err := New(cfg)
	if err != nil {
		// In example/test environment, just return instead of fatal
		return
	}
	defer db.Close()
	
	// Run migrations using the migrator
	migrator := migrations.NewMigrator(db.DB())
	if err := migrator.AutoMigrate(); err != nil {
		log.Printf("Failed to run migrations: %v", err)
		return
	}
	
	err = db.Ping(ctx)
	if err != nil {
		// In example/test environment, just return instead of fatal
		return
	}
}

// ExampleUsage demonstrates comprehensive database usage
func ExampleUsage() {
	ctx := context.Background()
	
	// Initialize configuration
	cfg := config.DatabaseConfig{
		Host:            "localhost",
		Port:            5432,
		User:            "postgres",
		Password:        "password",
		DBName:          "myapp",
		SSLMode:         "disable",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5,
		ConnMaxIdleTime: 1,
		LogLevel:        "info",
	}

	// Create database instance
	db, err := New(cfg)
	if err != nil {
		log.Fatal("Failed to create database:", err)
	}
	defer db.Close()

	// Check connection
	if err := db.Ping(ctx); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Run migrations using the migrator
	migrator := migrations.NewMigrator(db.DB())
	if err := migrator.AutoMigrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Use GORM directly
	gormDB := db.DB()
	
	// Create a user
	user := &models.User{
		Email:        "john@example.com",
		Username:     "johndoe",
		FirstName:    "John",
		LastName:     "Doe",
		PasswordHash: "hashed_password_here", // In real app, use bcrypt
	}
	if err := gormDB.Create(user).Error; err != nil {
		log.Fatal("Failed to create user:", err)
	}

	// Create a post with transaction
	err = db.Transaction(ctx, func(tx *gorm.DB) error {
		// Create post
		post := &models.Post{
			Title:     "My First Post",
			Slug:      "my-first-post",
			Content:   "This is the content of my first post.",
			Summary:   "A brief summary of the post",
			Published: true,
			UserID:    user.ID,
		}
		if err := tx.Create(post).Error; err != nil {
			return err
		}

		// Create tags
		tags := []models.Tag{
			{Name: "golang", Slug: "golang"},
			{Name: "tutorial", Slug: "tutorial"},
		}
		for _, tag := range tags {
			if err := tx.FirstOrCreate(&tag, models.Tag{Slug: tag.Slug}).Error; err != nil {
				return err
			}
		}

		// Associate tags with post
		return tx.Model(post).Association("Tags").Append(tags)
	})
	if err != nil {
		log.Fatal("Transaction failed:", err)
	}

	// Query examples with context
	dbWithCtx := db.WithContext(ctx)
	
	// Find all published posts with user and tags
	var posts []models.Post
	if err := dbWithCtx.DB().
		Preload("User").
		Preload("Tags").
		Where("published = ?", true).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		log.Fatal("Failed to query posts:", err)
	}

	// Example of using model methods
	if len(posts) > 0 {
		// Increment view count for the first post
		if err := posts[0].IncrementViewCount(gormDB); err != nil {
			log.Fatal("Failed to increment view count:", err)
		}
	}

	fmt.Printf("Found %d published posts\n", len(posts))
}