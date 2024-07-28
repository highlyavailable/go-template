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
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}
