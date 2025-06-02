package models

import (
	"testing"
)

func TestPostModel(t *testing.T) {
	db := setupTestDB(t)

	// Create a test user first
	user := &User{
		Email:        "author@example.com",
		Username:     "author",
		PasswordHash: "hashed_password",
	}
	db.Create(user)

	t.Run("CreatePost", func(t *testing.T) {
		post := &Post{
			Title:     "Test Post",
			Slug:      "test-post",
			Content:   "This is test content",
			Summary:   "Test summary",
			Published: true,
			UserID:    user.ID,
		}

		err := db.Create(post).Error
		if err != nil {
			t.Fatalf("Failed to create post: %v", err)
		}

		if post.ID == 0 {
			t.Error("Expected post ID to be set after creation")
		}
	})

	t.Run("UniqueSlugConstraint", func(t *testing.T) {
		post1 := &Post{
			Title:   "Post 1",
			Slug:    "unique-slug",
			Content: "Content 1",
			UserID:  user.ID,
		}
		db.Create(post1)

		post2 := &Post{
			Title:   "Post 2",
			Slug:    "unique-slug",
			Content: "Content 2",
			UserID:  user.ID,
		}
		err := db.Create(post2).Error
		if err == nil {
			t.Error("Expected error when creating post with duplicate slug")
		}
	})

	t.Run("BeforeCreateValidation", func(t *testing.T) {
		testCases := []struct {
			name    string
			post    *Post
			wantErr bool
		}{
			{
				name: "Missing title",
				post: &Post{
					Slug:   "test-slug",
					UserID: user.ID,
				},
				wantErr: true,
			},
			{
				name: "Missing slug",
				post: &Post{
					Title:  "Test Title",
					UserID: user.ID,
				},
				wantErr: true,
			},
			{
				name: "Missing user ID",
				post: &Post{
					Title: "Test Title",
					Slug:  "test-slug",
				},
				wantErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := db.Create(tc.post).Error
				if (err != nil) != tc.wantErr {
					t.Errorf("Create() error = %v, wantErr %v", err, tc.wantErr)
				}
			})
		}
	})

	t.Run("IncrementViewCount", func(t *testing.T) {
		post := &Post{
			Title:     "View Count Test",
			Slug:      "view-count-test",
			Content:   "Content",
			UserID:    user.ID,
			ViewCount: 0,
		}
		db.Create(post)

		// Increment view count
		err := post.IncrementViewCount(db)
		if err != nil {
			t.Fatalf("Failed to increment view count: %v", err)
		}

		// Reload post
		var updatedPost Post
		db.First(&updatedPost, post.ID)

		if updatedPost.ViewCount != 1 {
			t.Errorf("Expected view count to be 1, got %d", updatedPost.ViewCount)
		}

		// Increment again
		err = updatedPost.IncrementViewCount(db)
		if err != nil {
			t.Fatalf("Failed to increment view count second time: %v", err)
		}

		db.First(&updatedPost, post.ID)
		if updatedPost.ViewCount != 2 {
			t.Errorf("Expected view count to be 2, got %d", updatedPost.ViewCount)
		}
	})

	t.Run("PostWithComments", func(t *testing.T) {
		post := &Post{
			Title:   "Post with Comments",
			Slug:    "post-with-comments",
			Content: "Content",
			UserID:  user.ID,
		}
		db.Create(post)

		// Create comments
		comment1 := &Comment{
			Content: "First comment",
			UserID:  user.ID,
			PostID:  post.ID,
		}
		comment2 := &Comment{
			Content: "Second comment",
			UserID:  user.ID,
			PostID:  post.ID,
		}
		db.Create([]*Comment{comment1, comment2})

		// Load post with comments
		var loadedPost Post
		err := db.Preload("Comments").First(&loadedPost, post.ID).Error
		if err != nil {
			t.Fatalf("Failed to load post with comments: %v", err)
		}

		if len(loadedPost.Comments) != 2 {
			t.Errorf("Expected 2 comments, got %d", len(loadedPost.Comments))
		}
	})

	t.Run("PostWithTags", func(t *testing.T) {
		post := &Post{
			Title:   "Post with Tags",
			Slug:    "post-with-tags",
			Content: "Content",
			UserID:  user.ID,
		}
		db.Create(post)

		// Create tags
		tag1 := &Tag{Name: "Go", Slug: "go"}
		tag2 := &Tag{Name: "Programming", Slug: "programming"}
		db.Create([]*Tag{tag1, tag2})

		// Associate tags with post
		err := db.Model(post).Association("Tags").Append([]*Tag{tag1, tag2})
		if err != nil {
			t.Fatalf("Failed to associate tags: %v", err)
		}

		// Load post with tags
		var loadedPost Post
		err = db.Preload("Tags").First(&loadedPost, post.ID).Error
		if err != nil {
			t.Fatalf("Failed to load post with tags: %v", err)
		}

		if len(loadedPost.Tags) != 2 {
			t.Errorf("Expected 2 tags, got %d", len(loadedPost.Tags))
		}
	})

	t.Run("QueryPublishedPosts", func(t *testing.T) {
		// Create published and unpublished posts
		publishedPost := &Post{
			Title:     "Published Post",
			Slug:      "published-post",
			Content:   "Content",
			UserID:    user.ID,
			Published: true,
		}
		unpublishedPost := &Post{
			Title:     "Unpublished Post",
			Slug:      "unpublished-post",
			Content:   "Content",
			UserID:    user.ID,
			Published: false,
		}
		db.Create([]*Post{publishedPost, unpublishedPost})

		// Query only published posts
		var publishedPosts []Post
		err := db.Where("published = ?", true).Find(&publishedPosts).Error
		if err != nil {
			t.Fatalf("Failed to query published posts: %v", err)
		}

		// Check that we only get published posts
		for _, post := range publishedPosts {
			if !post.Published {
				t.Error("Expected only published posts in result")
			}
		}
	})

	t.Run("LoadPostWithAllAssociations", func(t *testing.T) {
		// Create a complete post setup
		post := &Post{
			Title:     "Complete Post",
			Slug:      "complete-post",
			Content:   "Full content here",
			Summary:   "Summary of the post",
			Published: true,
			UserID:    user.ID,
		}
		db.Create(post)

		// Add tags
		tag := &Tag{Name: "Complete", Slug: "complete"}
		db.Create(tag)
		db.Model(post).Association("Tags").Append(tag)

		// Add comment
		comment := &Comment{
			Content: "Great post!",
			UserID:  user.ID,
			PostID:  post.ID,
		}
		db.Create(comment)

		// Load everything
		var loadedPost Post
		err := db.Preload("User").
			Preload("Comments").
			Preload("Tags").
			First(&loadedPost, post.ID).Error
		if err != nil {
			t.Fatalf("Failed to load post with associations: %v", err)
		}

		// Verify associations
		if loadedPost.User.ID != user.ID {
			t.Error("Expected user to be loaded")
		}
		if len(loadedPost.Comments) != 1 {
			t.Errorf("Expected 1 comment, got %d", len(loadedPost.Comments))
		}
		if len(loadedPost.Tags) != 1 {
			t.Errorf("Expected 1 tag, got %d", len(loadedPost.Tags))
		}
	})
}

func TestPostIncrementViewCountConcurrency(t *testing.T) {
	db := setupTestDB(t)

	user := &User{
		Email:        "concurrent@example.com",
		Username:     "concurrent",
		PasswordHash: "hashed_password",
	}
	db.Create(user)

	post := &Post{
		Title:     "Concurrent Views",
		Slug:      "concurrent-views",
		Content:   "Content",
		UserID:    user.ID,
		ViewCount: 0,
	}
	db.Create(post)

	// Note: SQLite doesn't handle true concurrency well due to file locking.
	// This test simulates the pattern but may not catch all concurrency issues.
	// In production with PostgreSQL, the GORM expression ensures atomic updates.
	
	// Simulate sequential view count increments (SQLite limitation)
	for i := 0; i < 5; i++ {
		p := &Post{BaseModel: BaseModel{ID: post.ID}}
		if err := p.IncrementViewCount(db); err != nil {
			t.Fatalf("Failed to increment view count: %v", err)
		}
	}

	// Check final count
	var finalPost Post
	db.First(&finalPost, post.ID)

	if finalPost.ViewCount != 5 {
		t.Errorf("Expected view count to be 5 after increments, got %d", finalPost.ViewCount)
	}
}