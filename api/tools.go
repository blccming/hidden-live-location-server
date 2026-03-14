package api

import (
	crand "crypto/rand"
	"time"

	"github.com/blccming/hidden-live-location-server/db"
	"github.com/rs/zerolog/log"
	glide "github.com/valkey-io/valkey-glide/go/v2"
)

var startup_time = time.Now()

func getRuntime() string {
	return time.Since(startup_time).String()
}

/* tokens */
func tokenGenerate() (string, error) {
	const length = 6
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// rand bytes
	rb := make([]byte, length)
	if _, err := crand.Read(rb); err != nil {
		return "", err
	}

	// token in bytes
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[int(rb[i])%len(charset)]
	}

	return string(b), nil
}

func tokenCreate(g *glide.Client) string {
	for {
		token, err := tokenGenerate()
		if err != nil {
			log.Error().Stack().Err(err).Msg("Failed to generate token.")
			continue
		}
		exists, err := db.SessionExists(g, token)
		if err != nil {
			log.Error().Stack().Err(err).Msg("Failed to check session existence.")
			continue
		}
		if !exists {
			return token
		}
	}
}
