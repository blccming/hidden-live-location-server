package main

import (
	"os"

	"github.com/blccming/hidden-live-location-server/api"
	"github.com/blccming/hidden-live-location-server/configuration"
	"github.com/blccming/hidden-live-location-server/db"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := configuration.Configure()

	dbClient, err := db.Connect(cfg.DBHost, cfg.DBPort, cfg.DBPass)
	if err != nil {
		log.Panic().Stack().Err(err).Msg("Failed to connect and configure database. Exiting ..")
		os.Exit(1)
	}

	r := api.InitEndpoints(cfg.LogLevel == "DEBUG" || cfg.LogLevel == "TRACE", dbClient)
	log.Info().Msgf("Server starting on %s.", cfg.Host+":"+cfg.Port)
	r.Run(cfg.Host + ":" + cfg.Port)
}
