package main

import (
	"fmt"
	"goapp/api/routes"
	_ "goapp/docs" // Import generated docs
	"goapp/pkg/logging"
	"goapp/pkg/otel"
	"log"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

type Specification struct {
	AppName     string `envconfig:"APP_NAME" default:"goapp"`
	ProjectRoot string `envconfig:"PROJECT_ROOT" default:"/Users/PeterWBryant/Repos/go-template"`
	EnvPath     string `envconfig:"ENV_PATH"`
	LogDirPath  string `envconfig:"LOG_DIR_PATH"`
	CertDirPath string `envconfig:"CERT_DIR_PATH"`
	Env         string `envconfig:"ENV" default:"development"`
	Port        int    `envconfig:"PORT" default:"8080"`
}

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
	// Load env vars into Specification struct
	var s Specification
	err := envconfig.Process("GO_APP", &s)
	if err != nil {
		log.Fatal(err.Error())
	}
	format := "AppName: %s\nProjectRoot: %s\nEnvPath: %s\nLogDirPath: %s\nCertDirPath: %s\nEnv: %s\nPort: %d\n"
	_, err = fmt.Printf(format, s.AppName, s.ProjectRoot, s.EnvPath, s.LogDirPath, s.CertDirPath, s.Env, s.Port)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Init zapLogger configuration
	var loggerConf = logging.LoggerConfig{
		Environment:      s.Env,
		AppLogPath:       filepath.Join(s.LogDirPath, "app.log"),
		ErrLogPath:       filepath.Join(s.LogDirPath, "error.log"),
		WriteStdout:      true,
		EnableStackTrace: false,
		MaxSize:          1,
		MaxBackups:       5,
		MaxAge:           30,
		Compress:         true,
	}
	logging.InitLogger(loggerConf)
	// logging.TestRotation(1e4)     // Test log rotation by dumping 10k error msgs

	// Initialize OTel Tracer and Meter
	shutdownTracer := otel.InitTracer()
	defer shutdownTracer()
	shutdownMeter := otel.InitMeter()
	defer shutdownMeter()

	// Init custom counter
	counter := otel.InitCustomCounter("custom_counter")
	otel.UpdateCounter(counter, 100)

	// Initialize the router
	router := routes.SetupRouter() // Create all routes
	logging.Info("Server started", zap.Int("port", s.Port))
	router.Use(otelgin.Middleware("goapp")) // Add OpenTelemetry middleware

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
