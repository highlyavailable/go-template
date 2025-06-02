package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goapp/internal/container"
	"goapp/web/templates/pages"
)

// HomeHandler handles the home page
type HomeHandler struct {
	container *container.Container
}

// NewHomeHandler creates a new home handler
func NewHomeHandler(c *container.Container) *HomeHandler {
	return &HomeHandler{container: c}
}

// Index renders the home page
func (h *HomeHandler) Index(c *gin.Context) {
	component := pages.Home()
	
	c.Header("Content-Type", "text/html")
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		h.container.Logger.Error("Failed to render home page", zap.Error(err))
		c.String(http.StatusInternalServerError, "Failed to render page")
		return
	}
}