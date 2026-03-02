package main

import (
	"github.com/blccming/private-positioning-server/api"
	"github.com/blccming/private-positioning-server/configuration"
	"github.com/rs/zerolog/log"
)

func main() {
	ginCfg := configuration.Configure()

	r := api.InitEndpoints()
	log.Info().Msgf("Server starting on %s.", ginCfg)
	r.Run(ginCfg)
}
