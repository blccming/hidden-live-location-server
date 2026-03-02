package configuration

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
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

func initializeZerolog() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

	log.Debug().Msgf("Configured Zerolog.")
}

func Configure() string {
	// Initialize zerolog before first use
	initializeZerolog()

	c := config{
		port:     getEnvPort(DefaultPort),
		host:     getEnvHost(DefaultHost),
		logLevel: getEnvLogLevel(DefaultLogLevel),
	}

	// Set gin mode to release if log level is not DEBUG
	if c.logLevel != "DEBUG" && c.logLevel != "TRACE" {
		gin.SetMode(gin.ReleaseMode)
		log.Info().Msgf("Set gin mode to release due to log level: %s", c.logLevel)
	} else {
		log.Debug().Msgf("Left gin mode at debug due to log level: %s", c.logLevel)
	}

	// Configure zerolog
	zerolog.SetGlobalLevel(logLevels[c.logLevel])

	// Run the gin server
	return fmt.Sprintf("%s:%s", c.host, c.port)
}
