package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	BaseModel
	Email         string     `gorm:"uniqueIndex;not null" json:"email"`
	Username      string     `gorm:"uniqueIndex;not null" json:"username"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	PasswordHash  string     `gorm:"not null" json:"-"`
	Active        bool       `gorm:"default:true" json:"active"`
	EmailVerified bool       `gorm:"default:false" json:"email_verified"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	
	// Associations
	Posts    []Post    `gorm:"foreignKey:UserID" json:"posts,omitempty"`
	Comments []Comment `gorm:"foreignKey:UserID" json:"comments,omitempty"`
}

// BeforeCreate hook for User model
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.Username == "" {
		return errors.New("username is required")
	}
	if u.PasswordHash == "" {
		return errors.New("password is required")
	}
	return nil
}

// FullName returns the user's full name
func (u *User) FullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Username
	}
	return u.FirstName + " " + u.LastName
}