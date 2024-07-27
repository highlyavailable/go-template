package main

import (
	"fmt"
	"goapp/internal/config"
	"goapp/internal/logging"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		panic("ENV_PATH environment variable is not set")
	}
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Printf("Error loading .env file at %s: %v\n", envPath, err)
		panic(err)
	}

	// List of required environment variables
	requiredVars := []string{"ENV", "ENV_PATH", "LOG_DIR_PATH"}

	// Check if required environment variables are set
	err = config.CheckRequiredEnvVars(requiredVars)
	if err != nil {
		panic(err)
	}

	// Custom logger struct see gomodule/internal/logging/logging.go
	newLogger := logging.LoggerConfig{
		Environment:      os.Getenv("ENV"),
		EnableStackTrace: false,
		MaxSize:          1,
		MaxBackups:       5,
		MaxAge:           30,
		Compress:         true,
		AppLogPath:       fmt.Sprintf("%s/app.log", os.Getenv("LOG_DIR_PATH")),
		ErrLogPath:       fmt.Sprintf("%s/error.log", os.Getenv("LOG_DIR_PATH")),
	}

	logging.InitLogger(newLogger) // Initialize the logger
	// logging.TestRotation(1e4)     // Test log rotation by dumping 10,000 error msgs

	logging.Info("LOGFILE - Good to go")
	fmt.Println("STDOUT - Good to go")
}
