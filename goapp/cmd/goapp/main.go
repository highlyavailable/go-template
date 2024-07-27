package main

import (
	"fmt"
	"goapp/internal/logging"

	"github.com/joho/godotenv"
)

const (
	envPath = "../.env"
)

func main() {
	// Load the .env file
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Println("No .env file found")
	}

	// Custom logger struct see gomodule/internal/logging/logging.go
	newLogger := logging.LoggerConfig{
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
		AppLogPath: "./logs/app.log",
		ErrLogPath: "./logs/error.log",
	}

	logging.InitLogger(newLogger) // Initialize the logger
	logging.Info("Good to go")
	logging.Error("This is an error message")
}
