package configuration

import (
	"os"
	"strconv"

	"github.com/rs/zerolog"
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
		if level != zerolog.NoLevel {
			return logLevel
		}
	}
	return defaultValue
}

func getEnvPort(defaultValue string) string {
	if port := os.Getenv("PORT"); port != "" {
		if portInt, err := strconv.Atoi(port); err == nil && portInt >= 1024 && portInt <= 65535 {
			return port
		}
	}
	return defaultValue
}

func getEnvHost(defaultValue string) string {
	if host := os.Getenv("HOST"); host != "" {
		return host
	}
	return defaultValue
}
