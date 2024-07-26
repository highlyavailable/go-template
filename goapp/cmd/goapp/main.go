package main

import (
	"fmt"
	"goapp/internal/logging"
	"os"

	"github.com/joho/godotenv"
)

// Expects a .env file in the root of the project
// with the following content:
// ENV=development

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found")
	}

	// Get the environment
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	logging.InitLogger(env) // Initialize the logger
	// logging.Info("This is an info message")
	// logging.Error("This is an error message")
}
