package migrations

import (
	"fmt"

	"goapp/internal/models"
	"gorm.io/gorm"
)

// Migrator handles database migrations
type Migrator struct {
	db *gorm.DB
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

// AutoMigrate runs auto-migration for all models
func (m *Migrator) AutoMigrate() error {
	models := []interface{}{
		&models.User{},
		&models.Post{},
		&models.Comment{},
		&models.Tag{},
	}

	for _, model := range models {
		if err := m.db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
	}

	// Create indexes
	if err := m.createIndexes(); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

// createIndexes creates custom indexes for better performance
func (m *Migrator) createIndexes() error {
	// Add composite indexes
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_published_created_at ON posts(published, created_at DESC) WHERE deleted_at IS NULL").Error; err != nil {
		return err
	}

	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_comments_post_id_created_at ON comments(post_id, created_at) WHERE deleted_at IS NULL").Error; err != nil {
		return err
	}

	return nil
}

// DropAllTables drops all tables (use with caution!)
func (m *Migrator) DropAllTables() error {
	return m.db.Migrator().DropTable(
		&models.Tag{},
		&models.Comment{},
		&models.Post{},
		&models.User{},
		"post_tags", // many2many join table
	)
}