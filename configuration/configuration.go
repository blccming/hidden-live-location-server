package configuration

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	DefaultPort     = "8080"
	DefaultHost     = "0.0.0.0"
	DefaultLogLevel = "DEBUG"
)

func initializeZerolog(loglevel zerolog.Level) {
	zerolog.SetGlobalLevel(loglevel)

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006/01/02 - 15:04:05",
	})

	log.Debug().Msgf("Configured Zerolog.")
}

func Configure() string {
	// Initialize zerolog before first use
	loglevel := getEnvLogLevel(DefaultLogLevel)
	initializeZerolog(logLevels[loglevel])

	if loglevel != DefaultLogLevel {
		log.Info().Msgf("Changed default log level due to env var: %s", loglevel)
	} else {
		log.Debug().Msgf("Using default log level: %s", loglevel)
	}

	// Set gin mode to release if log level is not DEBUG
	if loglevel != "DEBUG" && loglevel != "TRACE" {
		gin.SetMode(gin.ReleaseMode)
		log.Info().Msgf("Set gin mode to release due to log level: %s", loglevel)
	} else {
		log.Debug().Msgf("Left gin mode at debug due to log level: %s", loglevel)
	}

	// Get port and host from environment variables or use defaults and return the address
	port := getEnvPort(DefaultPort)
	host := getEnvHost(DefaultHost)
	return fmt.Sprintf("%s:%s", host, port)
}
