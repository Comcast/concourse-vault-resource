package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

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

	rand.Seed(time.Now().UTC().UnixNano())
	version := request.Version.Version
	if len(request.Version.Version) <= 0 {
		version = fmt.Sprintf("%d", rand.Intn(100))
	}

	response := models.Response{
		Metadata: nil,
		Version: models.Version{
			Version: version,
		},
	}

	// first argument on stdin is the working directory
	vault, err := resource.New(os.Args[1], request, logger)
	if err != nil {
		logger.Fatal().Err(err).
			Msg("error creating resource client")
	}

	err = vault.In()
	if err != nil {
		logger.Fatal().Err(err).
			Msg("error running file for write")
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		logger.Fatal().Err(err).
			Msg("writing response")
	}
}
