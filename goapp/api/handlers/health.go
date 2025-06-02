package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheckHandler godoc
// @Summary Health check
// @Description Do a health check
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *Handler) HealthCheckHandler(c *gin.Context) {
	h.Logger.Info("Health check endpoint called")
	
	// Check database connection
	if h.Database != nil {
		if err := h.Database.Ping(c.Request.Context()); err != nil {
			h.Logger.Errorf("Database health check failed: %v", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "DOWN",
				"error":  "Database connection failed",
			})
			return
		}
	} else {
		h.Logger.Warn("Database not configured, skipping database health check")
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}
