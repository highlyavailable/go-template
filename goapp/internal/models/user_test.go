package models

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{}, &Tag{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestUserModel(t *testing.T) {
	db := setupTestDB(t)

	t.Run("CreateUser", func(t *testing.T) {
		user := &User{
			Email:        "test@example.com",
			Username:     "testuser",
			FirstName:    "Test",
			LastName:     "User",
			PasswordHash: "hashed_password",
		}

		err := db.Create(user).Error
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		if user.ID == 0 {
			t.Error("Expected user ID to be set after creation")
		}
	})

	t.Run("UniqueEmailConstraint", func(t *testing.T) {
		user1 := &User{
			Email:        "duplicate@example.com",
			Username:     "user1",
			PasswordHash: "hashed_password",
		}
		db.Create(user1)

		user2 := &User{
			Email:        "duplicate@example.com",
			Username:     "user2",
			PasswordHash: "hashed_password",
		}
		err := db.Create(user2).Error
		if err == nil {
			t.Error("Expected error when creating user with duplicate email")
		}
	})

	t.Run("UniqueUsernameConstraint", func(t *testing.T) {
		user1 := &User{
			Email:        "user3@example.com",
			Username:     "duplicateusername",
			PasswordHash: "hashed_password",
		}
		db.Create(user1)

		user2 := &User{
			Email:        "user4@example.com",
			Username:     "duplicateusername",
			PasswordHash: "hashed_password",
		}
		err := db.Create(user2).Error
		if err == nil {
			t.Error("Expected error when creating user with duplicate username")
		}
	})

	t.Run("BeforeCreateValidation", func(t *testing.T) {
		testCases := []struct {
			name     string
			user     *User
			wantErr  bool
		}{
			{
				name: "Missing email",
				user: &User{
					Username:     "testuser",
					PasswordHash: "hashed_password",
				},
				wantErr: true,
			},
			{
				name: "Missing username",
				user: &User{
					Email:        "test@example.com",
					PasswordHash: "hashed_password",
				},
				wantErr: true,
			},
			{
				name: "Missing password",
				user: &User{
					Email:    "test@example.com",
					Username: "testuser",
				},
				wantErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := db.Create(tc.user).Error
				if (err != nil) != tc.wantErr {
					t.Errorf("Create() error = %v, wantErr %v", err, tc.wantErr)
				}
			})
		}
	})

	t.Run("FullNameMethod", func(t *testing.T) {
		testCases := []struct {
			name         string
			user         User
			expectedName string
		}{
			{
				name: "With first and last name",
				user: User{
					FirstName: "John",
					LastName:  "Doe",
					Username:  "johndoe",
				},
				expectedName: "John Doe",
			},
			{
				name: "Without names",
				user: User{
					Username: "johndoe",
				},
				expectedName: "johndoe",
			},
			{
				name: "Only first name",
				user: User{
					FirstName: "John",
					Username:  "johndoe",
				},
				expectedName: "John ",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				if got := tc.user.FullName(); got != tc.expectedName {
					t.Errorf("FullName() = %v, want %v", got, tc.expectedName)
				}
			})
		}
	})

	t.Run("SoftDelete", func(t *testing.T) {
		user := &User{
			Email:        "softdelete@example.com",
			Username:     "softdeleteuser",
			PasswordHash: "hashed_password",
		}
		db.Create(user)

		// Soft delete the user
		db.Delete(user)

		// Try to find the user normally (should not find)
		var foundUser User
		result := db.First(&foundUser, user.ID)
		if result.Error == nil {
			t.Error("Expected user to not be found after soft delete")
		}

		// Find with unscoped (should find)
		var unscopedUser User
		result = db.Unscoped().First(&unscopedUser, user.ID)
		if result.Error != nil {
			t.Errorf("Expected to find user with unscoped query: %v", result.Error)
		}
		if unscopedUser.DeletedAt.Time.IsZero() {
			t.Error("Expected DeletedAt to be set")
		}
	})

	t.Run("UserWithAssociations", func(t *testing.T) {
		user := &User{
			Email:        "author@example.com",
			Username:     "author",
			PasswordHash: "hashed_password",
		}
		db.Create(user)

		// Create posts for the user
		post1 := &Post{
			Title:   "Post 1",
			Slug:    "post-1",
			Content: "Content 1",
			UserID:  user.ID,
		}
		post2 := &Post{
			Title:   "Post 2",
			Slug:    "post-2",
			Content: "Content 2",
			UserID:  user.ID,
		}
		db.Create([]*Post{post1, post2})

		// Load user with posts
		var loadedUser User
		err := db.Preload("Posts").First(&loadedUser, user.ID).Error
		if err != nil {
			t.Fatalf("Failed to load user with posts: %v", err)
		}

		if len(loadedUser.Posts) != 2 {
			t.Errorf("Expected 2 posts, got %d", len(loadedUser.Posts))
		}
	})

	t.Run("LastLoginAt", func(t *testing.T) {
		user := &User{
			Email:        "login@example.com",
			Username:     "loginuser",
			PasswordHash: "hashed_password",
		}
		db.Create(user)

		// Update last login
		now := time.Now()
		db.Model(user).Update("last_login_at", now)

		// Reload user
		var updatedUser User
		db.First(&updatedUser, user.ID)

		if updatedUser.LastLoginAt == nil {
			t.Error("Expected LastLoginAt to be set")
		} else if updatedUser.LastLoginAt.Unix() != now.Unix() {
			t.Errorf("Expected LastLoginAt to be %v, got %v", now, *updatedUser.LastLoginAt)
		}
	})
}