package models

import (
	"errors"

	"gorm.io/gorm"
)

// Comment represents a comment on a post
type Comment struct {
	BaseModel
	Content  string  `gorm:"type:text;not null" json:"content"`
	UserID   uint    `gorm:"not null;index" json:"user_id"`
	PostID   uint    `gorm:"not null;index" json:"post_id"`
	ParentID *uint   `gorm:"index" json:"parent_id,omitempty"`
	
	// Associations
	User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Post     Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`
	Parent   *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Replies  []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

// BeforeCreate hook for Comment model
func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.Content == "" {
		return errors.New("content is required")
	}
	if c.UserID == 0 {
		return errors.New("user_id is required")
	}
	if c.PostID == 0 {
		return errors.New("post_id is required")
	}
	return nil
}