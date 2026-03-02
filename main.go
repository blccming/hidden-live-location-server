package main

import (
	"github.com/blccming/private-positioning-server/api"
	"github.com/blccming/private-positioning-server/configuration"
	"github.com/rs/zerolog/log"
)

func main() {
	ginCfg, debugMode := configuration.Configure()

	r := api.InitEndpoints(debugMode) // Use swagger if debugMode (logLevel is DEBUG or TRACE)
	log.Info().Msgf("Server starting on %s.", ginCfg)
	r.Run(ginCfg)
}
