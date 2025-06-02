package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goapp/api/routes"
	_ "goapp/docs" // Import generated docs
	"goapp/internal/container"
	"goapp/internal/observability"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// @title GoApp REST API
// @version 1.0
// @description Production-ready Go REST API with dependency injection
// @termsOfService <url>

// @contact.name Peter Bryant
// @contact.url <url>
// @contact.email <email>
// @license.name Apache 2.0
// @license.url <url>

// @host localhost:8080
// @BasePath /
func main() {
	// Initialize dependency injection container
	c, err := container.New()
	if err != nil {
		fmt.Printf("Failed to initialize container: %v\n", err)
		os.Exit(1)
	}
	defer c.Close()

	c.Logger.Infof("Application starting: name=%s, env=%s", c.Config.App.Name, c.Config.App.Env)

	// Set Gin mode based on environment
	if c.Config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize OpenTelemetry if enabled
	var shutdownTracer, shutdownMeter func()
	if c.Config.Observability.Enabled {
		c.Logger.Info("OpenTelemetry enabled - initializing tracing and metrics")
		shutdownTracer = observability.InitTracer(c.Config.Observability)
		defer shutdownTracer()
		shutdownMeter = observability.InitMeter(c.Config.Observability)
		defer shutdownMeter()
		
		// Initialize custom counter for demonstration
		counter := observability.InitCustomCounter("http_requests_total")
		observability.UpdateCounter(counter, 1)
	}

	// Setup router with dependency injection
	router := routes.SetupRouter(c)
	
	// Add OpenTelemetry middleware if enabled
	if c.Config.Observability.Enabled {
		router.Use(otelgin.Middleware(c.Config.Observability.ServiceName))
	}

	// Setup HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.Config.App.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		c.Logger.Infof("Starting HTTP server on port %d", c.Config.App.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			c.Logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	c.Logger.Info("Shutting down server...")

	// Give outstanding requests 5 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		c.Logger.Fatalf("Server forced to shutdown: %v", err)
	}

	c.Logger.Info("Server exited")
}
