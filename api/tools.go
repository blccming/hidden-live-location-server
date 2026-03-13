package api

import (
	"math/rand"
	"os"
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
func tokenGenerate() string {
	const length = 6
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano() & int64(os.Getpid())))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func tokenCreate(g *glide.Client) string {
	for {
		token := tokenGenerate()
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
