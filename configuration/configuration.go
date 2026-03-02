package configuration

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// TODO: set zerolog level

const (
	DefaultPort     = "8080"
	DefaultHost     = "0.0.0.0"
	DefaultLogLevel = "DEBUG"
)

type config struct {
	port     string
	host     string
	logLevel string
}

func configureZerolog(level zerolog.Level) {
	zerolog.SetGlobalLevel(level)

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})
}

func Configure() string {
	c := config{
		port:     getEnvPort(DefaultPort),
		host:     getEnvHost(DefaultHost),
		logLevel: getEnvLogLevel(DefaultLogLevel),
	}

	// Set gin mode to release if log level is not DEBUG
	if c.logLevel != "DEBUG" && c.logLevel != "TRACE" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Configure zerolog
	configureZerolog(logLevels[c.logLevel])

	// Run the gin server
	return fmt.Sprintf("%s:%s", c.host, c.port)
}
