package configuration

import (
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logLevels = map[string]zerolog.Level{
	"TRACE": zerolog.TraceLevel,
	"DEBUG": zerolog.DebugLevel,
	"INFO":  zerolog.InfoLevel,
	"WARN":  zerolog.WarnLevel,
	"ERROR": zerolog.ErrorLevel,
	"FATAL": zerolog.FatalLevel,
	"PANIC": zerolog.PanicLevel,
	"NOLOG": zerolog.NoLevel,
}

func getEnvLogLevel(defaultValue string) string {
	logLevel := os.Getenv("LOGLEVEL")
	if level, exists := logLevels[logLevel]; exists {
		if level != logLevels[defaultValue] {
			return logLevel
		}
	}
	return defaultValue
}

func getEnvPort(defaultValue string) string {
	if port := os.Getenv("PORT"); port != "" {
		if portInt, err := strconv.Atoi(port); err == nil && portInt >= 1024 && portInt <= 65535 {
			log.Info().Msgf("Changed default port due to env var: %s", port)
			return port
		}
	}
	log.Info().Msgf("Using default port (might be due to port env var being out of range: 1024-65535): %s", defaultValue)
	return defaultValue
}

func getEnvHost(defaultValue string) string {
	if host := os.Getenv("HOST"); host != "" {
		log.Info().Msgf("Changed default host due to env var: %s", host)
		return host
	}
	log.Info().Msgf("Using default host: %s", defaultValue)
	return defaultValue
}
