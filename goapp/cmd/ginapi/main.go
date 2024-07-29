package main

import (
	"fmt"
	"goapp/api/routes"
	"goapp/internal/config"
	"goapp/internal/logging"
	"goapp/internal/otel"
	"log"
	"os"

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
	requiredVars := []string{"ENV", "ENV_PATH", "LOG_DIR_PATH", "CERT_DIR_PATH"}
	config.LoadEnv(requiredVars) // Panic if not found

	// Init zap logger
	newLogger := logging.LoggerConfig{
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
	logging.InitLogger(newLogger)
	// logging.TestRotation(1e4)     // Test log rotation by dumping 10k error msgs

	// Initialize OTel Tracer and Meter
	shutdownTracer := otel.InitTracer()
	defer shutdownTracer()
	shutdownMeter := otel.InitMeter()
	defer shutdownMeter()

	// Create and test new insecure http client
	// _, err := clients.NewInsecureClient("http://google.com")
	// if err != nil {
	// 	logging.Error("Error creating insecure client", zap.Error(err))
	// 	fmt.Println("Error creating insecure client:", err)
	// 	return
	// }

	// Create a new secure client
	// certDirPath := os.Getenv("CERT_DIR_PATH")
	// certPath := fmt.Sprintf("%s/cert.pem", certDirPath)
	// keyPath := fmt.Sprintf("%s/key.pem", certDirPath)
	// logging.Info("Cert path", zap.String("cert_path", certPath))
	// logging.Info("Key path", zap.String("key_path", keyPath))
	// secureClientConfig := clients.SecureClientConfig{
	// 	CertFile: certPath,
	// 	KeyFile:  keyPath,
	// 	// ProxyURL:       "http://proxy.example.com",
	// 	URLForConnTest: "http://google.com",
	// }
	// _, err = clients.NewSecureClient(secureClientConfig)
	// if err != nil {
	// 	logging.Error("Error creating secure client", zap.Error(err))
	// 	fmt.Println("Error creating secure client:", err)
	// 	return
	// }

	// Initialize the router
	router := routes.SetupRouter() // Create all routes
	logging.Info("Server started", zap.String("port", "8080"))
	router.Use(otelgin.Middleware("goapp")) // Add OpenTelemetry middleware

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
