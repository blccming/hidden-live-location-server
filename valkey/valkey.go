package valkey

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/valkey-io/valkey-go"
)

func InitializeValkey() {
	// TODO: make these configurable via environment variable
	address := "localhost:6379"
	password := "YOUR_PASSWORD_HERE"

	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{address}, Password: password})
	if err != nil {
		log.Panic().
			Err(err).
			Msg("failed to initialize valkey client")
	}
	defer client.Close()

	ctx := context.Background()

	serverConfig(client, ctx)

	// value setting/getting tests
	//
	// err = client.Do(ctx, client.B().Set().Key("abc").Value("123").Nx().Build()).Error()
	// err = client.Do(ctx, client.B().Set().Key("def").Value("456").Nx().Build()).Error()

	strsl, err := client.Do(ctx, client.B().Get().Key("def").Build()).ToString()
	if err != nil {
		log.Error().Err(err).Msg("failed to get abc")
	}
	fmt.Print(strsl)
}

func serverConfig(client valkey.Client, ctx context.Context) {
	settings := map[string]string{
		"maxmemory":  "4GB", // TODO: make configurable via environment variable
		"save":       "",
		"appendonly": "no",
	}

	for setting, value := range settings {
		if err := client.Do(ctx, client.B().ConfigSet().ParameterValue().ParameterValue(setting, value).Build()).Error(); err != nil {
			log.Panic().Err(err).Msgf("failed to set %s", setting)
		}
	}
}
