package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goapp/internal/container"
	"goapp/internal/models"
	"goapp/web/templates/pages"
)

// PostsHandler handles post-related web pages
type PostsHandler struct {
	container *container.Container
}

// NewPostsHandler creates a new posts handler
func NewPostsHandler(c *container.Container) *PostsHandler {
	return &PostsHandler{container: c}
}

// Index renders the posts list page
func (h *PostsHandler) Index(c *gin.Context) {
	var posts []models.Post
	
	// Fetch posts from database
	if h.container.Database != nil {
		db := h.container.Database.DB()
		if err := db.Preload("User").Order("created_at DESC").Find(&posts).Error; err != nil {
			h.container.Logger.Error("Failed to fetch posts", zap.Error(err))
			c.String(http.StatusInternalServerError, "Failed to fetch posts")
			return
		}
	} else {
		// Mock data when database is not available
		h.container.Logger.Warn("Database not available, using mock data")
		posts = []models.Post{
			{
				BaseModel: models.BaseModel{ID: 1},
				Title:     "Welcome to GoApp",
				Slug:      "welcome-to-goapp",
				Summary:   "This is a demo post showing the web UI capabilities",
				Published: true,
				ViewCount: 42,
			},
			{
				BaseModel: models.BaseModel{ID: 2},
				Title:     "Building with HTMX and Templ",
				Slug:      "building-with-htmx-templ",
				Summary:   "Learn how to build dynamic web apps with Go",
				Published: true,
				ViewCount: 128,
			},
		}
	}
	
	component := pages.PostsIndex(posts)
	
	c.Header("Content-Type", "text/html")
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		h.container.Logger.Error("Failed to render posts page", zap.Error(err))
		c.String(http.StatusInternalServerError, "Failed to render page")
		return
	}
}