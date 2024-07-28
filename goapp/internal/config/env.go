package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// loadEnv loads the environment variables from the specified .env file path and checks if the required variables are set.
// Requires the ENV_PATH environment variable to be set, takes a list of required environment variables.
func LoadEnv(requiredVars []string) {
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

	err = checkRequiredEnvVars(requiredVars)
	if err != nil {
		panic(err)
	}
}

// CheckRequiredEnvVars checks if the required environment variables are set.
// It takes a slice of requiredVars, which contains the names of the environment variables to check.
// If any of the required environment variables are not set, it returns an error indicating the missing variables.
// Otherwise, it returns nil.
func checkRequiredEnvVars(requiredVars []string) error {
	var missingVars []string
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			missingVars = append(missingVars, v)
		}
		fmt.Println(v, " = ", os.Getenv(v))
	}
	if len(missingVars) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missingVars, ", "))
	}
	return nil
}
