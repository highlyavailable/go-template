package routes

import (
	"goapp/api/handlers"
	"goapp/api/handlers/web"
	"goapp/internal/container"

	_ "goapp/docs" // Import generated docs

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger"

	swaggerFiles "github.com/swaggo/files"
)

// SetupRouter sets up the Gin router with all the routes
func SetupRouter(container *container.Container) *gin.Engine {
	router := gin.Default()
	
	// Static files
	router.Static("/static", "./web/static")
	
	// Initialize handlers with dependency injection
	h := handlers.New(container)
	
	// Initialize web handlers
	homeHandler := web.NewHomeHandler(container)
	postsHandler := web.NewPostsHandler(container)
	partialsHandler := web.NewPartialsHandler(container)
	
	// API routes
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // Swagger endpoint
	container.Logger.Info("Swagger docs available at http://localhost:8080/swagger/index.html")

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check endpoint
	router.GET("/health", h.HealthCheckHandler)
	
	// Web routes
	router.GET("/", homeHandler.Index)
	router.GET("/posts", postsHandler.Index)
	
	// Partial routes for HTMX
	partials := router.Group("/partials")
	{
		partials.GET("/activity-feed", partialsHandler.ActivityFeed)
		partials.GET("/notifications", partialsHandler.Notifications)
		partials.GET("/user-menu", partialsHandler.UserMenu)
		partials.POST("/notifications/:id/read", partialsHandler.MarkNotificationRead)
	}

	return router
}
