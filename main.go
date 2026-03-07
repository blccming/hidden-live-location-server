package main

import (
	"github.com/blccming/hidden-live-location-server/api"
	"github.com/blccming/hidden-live-location-server/configuration"
	"github.com/rs/zerolog/log"
)

func main() {
	ginCfg, debugMode := configuration.Configure()

	r := api.InitEndpoints(debugMode) // Use swagger if debugMode (logLevel is DEBUG or TRACE)
	log.Info().Msgf("Server starting on %s.", ginCfg)
	r.Run(ginCfg)
}
