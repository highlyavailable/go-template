package routes

import (
	"goapp/api/handlers"
	"goapp/pkg/logging"

	_ "goapp/docs" // Import generated docs

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger"

	swaggerFiles "github.com/swaggo/files"
)

// SetupRouter sets up the Gin router with all the routes
func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // Swagger endpoint
	logging.Info("Swagger docs available at http://localhost:8080/swagger/index.html")

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check endpoint
	router.GET("/health", handlers.HealthCheckHandler)

	return router
}
