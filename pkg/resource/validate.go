package resource

import (
	"errors"
	"os"
	"strings"

	"github.com/comcast/concourse-vault-resource/pkg/resource/models"
)

// validate - validates the resource configuration
func validate(config models.Request) (models.Request, error) {
	if len(config.Source.VaultAddr) <= 0 {
		config.Source.VaultAddr = os.Getenv("VAULT_ADDR")
		if len(config.Source.VaultAddr) <= 0 {
			return config, errors.New("required argument vault_addr was not provided")
		}
	}

	if len(config.Source.VaultPaths) <= 0 {
		return config, errors.New("required argument vault_paths was not provided")
	}

	if len(config.Source.Format) <= 0 {
		config.Source.Format = "json"
	}

	if config.Source.Retries <= 0 {
		config.Source.Retries = 3
	}

	if !strings.Contains(config.Source.Format, "json") &&
		!strings.Contains(config.Source.Format, "yaml") {
		return config, errors.New("format provided is not supported. supported output formats are : \"json\" or \"yaml\"")
	}

	if len(config.Source.VaultToken) <= 0 && len(config.Source.SecretID) <= 0 {
		config.Source.VaultToken = os.Getenv("VAULT_TOKEN")
		if len(config.Source.VaultToken) <= 0 {
			return config, errors.New("vault_token was not provided")
		}
	}

	return config, nil
}
