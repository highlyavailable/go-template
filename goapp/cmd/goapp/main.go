package main

import (
	"fmt"
	"goapp/internal/config"
	"goapp/internal/logging"
	"goapp/internal/otel"
	"os"

	_ "goapp/docs" // Import generated docs
)

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

	logging.Info("GoApp is running")
}
