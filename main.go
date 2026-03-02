package main

import (
	"github.com/blccming/private-positioning-server/api"
	"github.com/blccming/private-positioning-server/configuration"
)

func main() {
	ginCfg := configuration.Configure()

	r := api.InitEndpoints()
	r.Run(ginCfg)
}
