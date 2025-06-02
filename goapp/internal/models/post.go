package models

import (
	"errors"

	"gorm.io/gorm"
)

// Post represents a blog post or article
type Post struct {
	BaseModel
	Title       string `gorm:"not null" json:"title"`
	Slug        string `gorm:"uniqueIndex;not null" json:"slug"`
	Content     string `gorm:"type:text" json:"content"`
	Summary     string `gorm:"type:text" json:"summary"`
	Published   bool   `gorm:"default:false;index" json:"published"`
	ViewCount   uint   `gorm:"default:0" json:"view_count"`
	UserID      uint   `gorm:"not null;index" json:"user_id"`
	
	// Associations
	User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	Tags     []Tag     `gorm:"many2many:post_tags;" json:"tags,omitempty"`
}

// BeforeCreate hook for Post model
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.Title == "" {
		return errors.New("title is required")
	}
	if p.Slug == "" {
		return errors.New("slug is required")
	}
	if p.UserID == 0 {
		return errors.New("user_id is required")
	}
	return nil
}

// IncrementViewCount increments the view count for the post
func (p *Post) IncrementViewCount(db *gorm.DB) error {
	return db.Model(p).Update("view_count", gorm.Expr("view_count + ?", 1)).Error
}