package models

// Tag represents a tag that can be assigned to posts
type Tag struct {
	BaseModel
	Name  string `gorm:"uniqueIndex;not null" json:"name"`
	Slug  string `gorm:"uniqueIndex;not null" json:"slug"`
	
	// Associations
	Posts []Post `gorm:"many2many:post_tags;" json:"posts,omitempty"`
}