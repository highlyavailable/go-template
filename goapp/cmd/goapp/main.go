package main

import (
	_ "goapp/docs" // Import generated docs
	"goapp/pkg/logging"
	"log"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
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

var s Specification

func init() {
	// Load env vars into Specification struct
	err := envconfig.Process("GO_APP", &s)
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
}

func main() {
	logging.Info("Hello, World!")
}
