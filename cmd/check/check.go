package main

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog"

	"github.com/comcast/concourse-vault-resource/pkg/resource"
	"github.com/comcast/concourse-vault-resource/pkg/resource/models"
)

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = ""
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	var request models.Request
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		logger.Fatal().Err(err).
			Msg("error reading from stdin")
	}

	vault, err := resource.New("vault", request, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("error creating resource client")
	}

	if err := json.NewEncoder(os.Stdout).Encode(vault.Check()); err != nil {
		logger.Fatal().Err(err).
			Msg("writing response")
	}
}
