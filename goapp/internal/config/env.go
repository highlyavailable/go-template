package config

import (
	"fmt"
	"os"
	"strings"
)

// CheckRequiredEnvVars checks if the required environment variables are set.
// It takes a slice of requiredVars, which contains the names of the environment variables to check.
// If any of the required environment variables are not set, it returns an error indicating the missing variables.
// Otherwise, it returns nil.
func CheckRequiredEnvVars(requiredVars []string) error {
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
