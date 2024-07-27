package main

import (
	"fmt"
	"goapp/internal/clients"
	"goapp/internal/config"
	"goapp/internal/logging"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
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
	requiredVars := []string{"ENV", "ENV_PATH", "LOG_DIR_PATH", "CERT_DIR_PATH"}

	// Check if required environment variables are set
	err = config.CheckRequiredEnvVars(requiredVars)
	if err != nil {
		panic(err)
	}

	// Custom logger struct see gomodule/internal/logging/logging.go
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

	logging.InitLogger(newLogger) // Initialize the logger
	// logging.TestRotation(1e4)     // Test log rotation by dumping 10,000 error msgs

	// Create a new insecure client (does not verify TLS certs)
	client, err := clients.NewInsecureClient()
	if err != nil {
		logging.Error("Error creating insecure client", zap.Error(err))
		fmt.Println("Error creating insecure client:", err)
		return
	}

	fmt.Println(client.CheckRedirect)

	// Create a new secure client
	// certDirPath := os.Getenv("CERT_DIR_PATH")
	// certPath := fmt.Sprintf("%s/cert.pem", certDirPath)
	// keyPath := fmt.Sprintf("%s/key.pem", certDirPath)
	// logging.Info("Cert path", zap.String("cert_path", certPath))
	// logging.Info("Key path", zap.String("key_path", keyPath))

	// client, err := clients.NewSecureClient("http://proxy.example.com", certPath, keyPath)
	// if err != nil {
	// 	logging.Error("Error creating secure client", zap.Error(err))
	// 	fmt.Println("Error creating secure client:", err)
	// 	return
	// }

}
