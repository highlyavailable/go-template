package web

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goapp/internal/container"
	"goapp/web/templates/partials"
)

// PartialsHandler handles partial template requests
type PartialsHandler struct {
	container *container.Container
}

// NewPartialsHandler creates a new partials handler
func NewPartialsHandler(c *container.Container) *PartialsHandler {
	return &PartialsHandler{container: c}
}

// ActivityFeed renders the activity feed partial
func (h *PartialsHandler) ActivityFeed(c *gin.Context) {
	// Mock activity data - in real app, fetch from database
	activities := []partials.ActivityItem{
		{
			Type:        "post_created",
			Description: "John Doe created a new post \"Getting Started with Go\"",
			Time:        time.Now().Add(-10 * time.Minute),
			Icon:        "M12 6v6m0 0v6m0-6h6m-6 0H6",
			IconColor:   "bg-blue-500",
		},
		{
			Type:        "comment_added",
			Description: "Jane Smith commented on \"Building REST APIs\"",
			Time:        time.Now().Add(-1 * time.Hour),
			Icon:        "M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z",
			IconColor:   "bg-green-500",
		},
		{
			Type:        "user_joined",
			Description: "New user Alice Johnson joined the platform",
			Time:        time.Now().Add(-3 * time.Hour),
			Icon:        "M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z",
			IconColor:   "bg-purple-500",
		},
	}
	
	component := partials.ActivityFeed(activities)
	
	c.Header("Content-Type", "text/html")
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		h.container.Logger.Error("Failed to render activity feed", zap.Error(err))
		c.String(http.StatusInternalServerError, "Failed to render partial")
		return
	}
}

// Notifications renders the notifications partial
func (h *PartialsHandler) Notifications(c *gin.Context) {
	// Mock notification data - in real app, fetch from database
	notifications := []partials.Notification{
		{
			ID:      "1",
			Title:   "New comment on your post",
			Message: "Someone commented on \"Getting Started with Go\"",
			Type:    "info",
			Read:    false,
		},
		{
			ID:      "2",
			Title:   "Post published successfully",
			Message: "Your post is now live",
			Type:    "success",
			Read:    true,
		},
	}
	
	component := partials.NotificationsList(notifications)
	
	c.Header("Content-Type", "text/html")
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		h.container.Logger.Error("Failed to render notifications", zap.Error(err))
		c.String(http.StatusInternalServerError, "Failed to render partial")
		return
	}
}

// UserMenu renders the user menu dropdown
func (h *PartialsHandler) UserMenu(c *gin.Context) {
	// In real app, get username from session/context
	username := "John Doe"
	
	component := partials.UserMenuDropdown(username)
	
	c.Header("Content-Type", "text/html")
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		h.container.Logger.Error("Failed to render user menu", zap.Error(err))
		c.String(http.StatusInternalServerError, "Failed to render partial")
		return
	}
}

// MarkNotificationRead handles marking a notification as read
func (h *PartialsHandler) MarkNotificationRead(c *gin.Context) {
	notifID := c.Param("id")
	
	// In real app, update notification status in database
	h.container.Logger.Info("Marking notification as read", zap.String("id", notifID))
	
	// Return empty response - the notification will be re-rendered
	c.Status(http.StatusOK)
}