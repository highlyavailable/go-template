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

	logging.InitLogger() // Initialize the logger
	logging.Info("Good to go")
	// logging.Error("This is an error message")
}
