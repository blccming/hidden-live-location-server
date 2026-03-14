package configuration

import (
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
	DefaultDBHost   = "localhost"
	DefaultDBPort   = 6379
	DefaultDBPass   = "changeme"
)

type AppConfig struct {
	LogLevel string `env:"LOG_LEVEL" default:"DEBUG"`
	Host     string `env:"HOST" default:"0.0.0.0"`
	Port     string `env:"PORT" default:"8080"`
	DBHost   string `env:"DB_HOST" default:"localhost"`
	DBPort   int    `env:"DB_PORT" default:"6379"`
	DBPass   string `env:"DB_PASS" default:"changeme"`
}

func initializeZerolog(loglevel zerolog.Level) {
	zerolog.SetGlobalLevel(loglevel)

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006/01/02 - 15:04:05",
	})

	log.Debug().Msgf("Configured Zerolog.")
}

func Configure() AppConfig {
	var cfg AppConfig

	// Initialize zerolog before first use
	cfg.LogLevel = getEnvLogLevel(DefaultLogLevel)
	initializeZerolog(logLevels[cfg.LogLevel])

	if cfg.LogLevel != DefaultLogLevel {
		log.Info().Msgf("Changed default log level due to env var: %s", cfg.LogLevel)
	} else {
		log.Debug().Msgf("Using default log level: %s", cfg.LogLevel)
	}

	// Set gin mode to release if log level is not DEBUG
	if cfg.LogLevel != "DEBUG" && cfg.LogLevel != "TRACE" {
		gin.SetMode(gin.ReleaseMode)
		log.Info().Msgf("Set gin mode to release due to log level: %s", cfg.LogLevel)
	} else {
		log.Debug().Msgf("Left gin mode at debug due to log level: %s", cfg.LogLevel)
	}

	// Get port and host from environment variables or use defaults and return the address
	cfg.Port = getEnvPort(DefaultPort)
	cfg.Host = getEnvHost(DefaultHost)
	cfg.DBHost = getEnvDBHost(DefaultDBHost)
	cfg.DBPort = getEnvDBPort(DefaultDBPort)
	cfg.DBPass = getEnvDBPass(DefaultDBPass)
	return cfg
}
