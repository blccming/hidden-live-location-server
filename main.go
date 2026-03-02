package main

import (
	"github.com/blccming/private-positioning-server/api"
)

func main() {
	r := api.InitEndpoints()
	r.Run("localhost:8080")
}
