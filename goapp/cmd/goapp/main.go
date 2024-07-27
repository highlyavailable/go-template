package main

import (
	"fmt"
	"goapp/internal/logging"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		fmt.Println("ENV_PATH environment variable is not set")
		panic("ENV_PATH environment variable is not set")
	}
	fmt.Println("ENV_PATH set to: ", envPath)
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Printf("Error loading .env file at %s: %v\n", envPath, err)
		panic(err)
	}

	// Custom logger struct see gomodule/internal/logging/logging.go
	newLogger := logging.LoggerConfig{
		Environment:      "development",
		EnableStackTrace: true,
		MaxSize:          1,
		MaxBackups:       5,
		MaxAge:           30,
		Compress:         true,
		AppLogPath:       "../logs/app.log",
		ErrLogPath:       "../logs/error.log",
	}

	logging.InitLogger(newLogger) // Initialize the logger
	logging.TestRotation(1e4)     // Test log rotation by dumping 10,000 error msgs

	logging.Info("Good to go")
	fmt.Println("Good to go")
}
