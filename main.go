package main

import (
	"os"

	"github.com/blccming/hidden-live-location-server/api"
	"github.com/blccming/hidden-live-location-server/configuration"
	"github.com/blccming/hidden-live-location-server/db"
	"github.com/rs/zerolog/log"
)

func main() {
	configuration.Configure() // remove later

	dbClient, err := db.Connect()
	if err != nil {
		log.Panic().Stack().Err(err).Msg("Failed to connect and configure database. Exiting ..")
		os.Exit(1)
	}

	// TESTING ONLY
	db.Test(dbClient)

	// stop them from executing for testing valkey atm
	if false {
		ginCfg, debugMode := configuration.Configure()

		r := api.InitEndpoints(debugMode) // Use swagger if debugMode (logLevel is DEBUG or TRACE)
		log.Info().Msgf("Server starting on %s.", ginCfg)
		r.Run(ginCfg)
	}
}
