package main

import (
	"fmt"
	"goapp/api/routes"
	"goapp/pkg/env"
	"goapp/pkg/logging"
	custotel "goapp/pkg/otel"
	"log"
	"os"

	// Import global package

	_ "goapp/docs" // Import generated docs

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

// Import generated docs

// @title GoApp Gin Rest API
// @version 1.0
// @description This is a sample server GoApp server.
// @termsOfService <url>

// @contact.name Peter Bryant
// @contact.url <url>
// @contact.email <email>
// @license.name Apache 2.0
// @license.url <url>

// @host localhost:8080
// @BasePath /
func main() {
	// List of required environment variables
	requiredVars := []string{"ENV_PATH", "ENV", "LOG_DIR_PATH", "CERT_DIR_PATH"}
	env.LoadEnv(requiredVars) // Panic if not found

	// Init zap logger
	loggerConf := logging.LoggerConfig{
		Environment:      os.Getenv("ENV"),
		WriteStdout:      true,
		EnableStackTrace: false,
		MaxSize:          1,
		MaxBackups:       5,
		MaxAge:           30,
		Compress:         true,
		AppLogPath:       fmt.Sprintf("%s/app.log", os.Getenv("LOG_DIR_PATH")),
		ErrLogPath:       fmt.Sprintf("%s/error.log", os.Getenv("LOG_DIR_PATH")),
	}
	logging.InitLogger(loggerConf)
	// logging.TestRotation(1e4)     // Test log rotation by dumping 10k error msgs

	// Initialize OTel Tracer and Meter
	shutdownTracer := custotel.InitTracer()
	defer shutdownTracer()
	shutdownMeter := custotel.InitMeter()
	defer shutdownMeter()

	// Init custom counter
	counter := custotel.InitCustomCounter("custom_counter")
	custotel.UpdateCounter(counter, 100)

	// Initialize the router
	router := routes.SetupRouter() // Create all routes
	logging.Info("Server started", zap.String("port", "8080"))
	router.Use(otelgin.Middleware("goapp")) // Add OpenTelemetry middleware

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
