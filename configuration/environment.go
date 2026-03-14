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

func getEnvHost(defaultValue string) string {
	if host := os.Getenv("HOST"); host != "" {
		log.Debug().Msgf("Changed default host due to env var: %s", host)
		return host
	}
	log.Debug().Msgf("Using default host: %s", defaultValue)
	return defaultValue
}

func getEnvPort(defaultValue string) string {
	if port := os.Getenv("PORT"); port != "" {
		if portInt, err := strconv.Atoi(port); err == nil && portInt >= 1024 && portInt <= 65535 {
			log.Debug().Msgf("Changed default port due to env var: %s", port)
			return port
		}
	}
	log.Debug().Msgf("Using default port (might be due to port env var being out of range: 1024-65535): %s", defaultValue)
	return defaultValue
}

func getEnvDBHost(defaultValue string) string {
	if host := os.Getenv("DB_HOST"); host != "" {
		log.Debug().Msgf("Changed default DB host due to env var: %s", host)
		return host
	}
	log.Debug().Msgf("Using default DB host: %s", defaultValue)
	return defaultValue
}

func getEnvDBPort(defaultValue int) int {
	if port := os.Getenv("DB_PORT"); port != "" {
		if portInt, err := strconv.Atoi(port); err == nil && portInt >= 1024 && portInt <= 65535 {
			log.Debug().Msgf("Changed default DB port due to env var: %s", port)
			return portInt
		}
	}
	log.Debug().Msgf("Using default DB port (might be due to port env var being out of range: 1024-65535): %d", defaultValue)
	return defaultValue
}

func getEnvDBPass(defaultValue string) string {
	if pass := os.Getenv("DB_PASSWORD"); pass != "" {
		log.Debug().Msgf("Changed default DB password due to env var: %s", pass)
		return pass
	}
	log.Debug().Msgf("Using default DB password: %s", defaultValue)
	return defaultValue
}
