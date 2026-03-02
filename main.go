package main

import (
	"github.com/blccming/private-positioning-server/api"
	"github.com/blccming/private-positioning-server/configuration"
	"github.com/blccming/private-positioning-server/valkey"
	"github.com/rs/zerolog/log"
)

func main() {
	configuration.Configure() // remove later

	valkey.InitializeValkey()

	// stop them from executing for testing valkey atm
	if false {
		ginCfg, debugMode := configuration.Configure()

		r := api.InitEndpoints(debugMode) // Use swagger if debugMode (logLevel is DEBUG or TRACE)
		log.Info().Msgf("Server starting on %s.", ginCfg)
		r.Run(ginCfg)
	}
}
