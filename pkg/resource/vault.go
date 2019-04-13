package resource

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/rs/zerolog"
	yaml "gopkg.in/yaml.v2"

	"github.com/comcast/concourse-vault-resource/pkg/resource/models"
)

// Vault - the vault resource interface
type Vault interface {
	Check() []models.Version
	In() error
}

// Resource - the vault resource
type Resource struct {
	client   *api.Client
	logger   zerolog.Logger
	config   models.Request
	secrets  map[string]interface{}
	workDir  string
	roleID   string
	secretID string
}

// New - returns a vault client for interaction with the vault API
func New(
	workDir string,
	config models.Request,
	logger zerolog.Logger,
) (*Resource, error) {
	var err error

	if config.Source.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		logger = logger.With().Caller().Logger()
	}

	config, err = validate(config)
	if err != nil {
		logger.Fatal().Err(err).
			Msg("error validating resource configuration")
	}

	c, err := api.NewClient(
		&api.Config{
			Address: config.Source.VaultAddr,
			HttpClient: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: config.Source.VaultInsecure,
					},
				},
			},
			MaxRetries: config.Source.Retries,
		},
	)
	if err != nil {
		logger.Fatal().AnErr("err", err).
			Msg("error occured creating client")
	}

	if len(config.Source.VaultToken) > 0 {
		c.SetToken(config.Source.VaultToken)
	}

	r := &Resource{
		client:  c,
		config:  config,
		logger:  logger,
		workDir: workDir,
		secrets: make(map[string]interface{}, 0),
	}

	err = r.setToken()
	if err != nil {
		r.logger.Fatal().Err(err).
			Msg("error logging into vault")
	}

	return r, nil
}

// Check - checks vault for new secret version
func (r Resource) Check() []models.Version {
	err := r.renewToken()
	if err != nil {
		r.logger.Fatal().AnErr("err", err).
			Msg("error occured renewing token")
	}

	var versions []models.Version
	for p, ver := range r.config.Source.VaultPaths {
		if ver > 0 || ver == -1 {
			s, err := r.client.Logical().Read(
				strings.Replace(p, "data", "metadata", 1),
			)
			if err != nil {
				r.logger.Fatal().AnErr("err", err).
					Msg("error occured reading paths")
			}

			if s != nil {
				versions = append(versions, models.Version{
					Path: p,
					Version: fmt.Sprintf(
						"%v", s.Data["current_version"],
					),
				})
			}
		} else {
			versions = append(versions, models.Version{
				Path:    p,
				Version: "1",
			})
		}
	}

	return versions
}

// In - executes the resource
func (r *Resource) In() error {
	err := r.renewToken()
	if err != nil {
		r.logger.Fatal().AnErr("err", err).
			Msg("error occured renewing token")
	}

	err = r.read()
	if err != nil {
		r.logger.Fatal().Err(err).
			Msg("error reading secrets")
	}

	r.prefix()

	r.sanitize()

	r.upcase()

	err = r.format()
	if err != nil {
		r.logger.Fatal().Err(err).
			Msg("error formatting secrets")
	}

	return nil
}

// format - formats the output in either json or yaml
func (r Resource) format() error {
	var (
		b   []byte
		err error
	)

	switch f := strings.ToLower(r.config.Source.Format); f {
	case "json":
		b, err = json.Marshal(r.secrets)
		if err != nil {
			return err
		}

	case "yaml":
		b, err = yaml.Marshal(r.secrets)
		if err != nil {
			return err
		}

	default:
		b, err = json.Marshal(r.secrets)
		if err != nil {
			return err
		}
	}

	if len(b) <= 0 {
		return errors.New("no secrets found to write to file")
	}

	err = r.write(b)
	if err != nil {
		return err
	}

	return nil
}

// prefix - adds a custom prefix to each key
func (r *Resource) prefix() {
	if len(r.config.Source.Prefix) <= 0 {
		return
	}

	s := make(map[string]interface{}, 0)
	for k, v := range r.secrets {
		s[fmt.Sprintf("%s_%s", r.config.Source.Prefix, k)] = v
	}
	r.secrets = s
}

// read - reads vault for a secret at a given path
func (r *Resource) read() error {
	var (
		s      *api.Secret
		err    error
		result = make(map[string]interface{}, 0)
	)

	for p, ver := range r.config.Source.VaultPaths {
		if ver > 0 {
			s, err = r.client.Logical().ReadWithData(p, map[string][]string{
				"version": []string{fmt.Sprintf("%d", ver)},
			})
		} else {
			s, err = r.client.Logical().Read(p)
		}
		if err != nil {
			r.logger.Fatal().AnErr("err", err).
				Msg("error occured reading paths")
		}

		if s != nil {
			// KV2
			if d, ok := s.Data["data"]; ok {
				switch t := d.(type) {
				case map[string]interface{}:
					for k, v := range t {
						result[k] = v
					}
				default:
					r.logger.Debug().Msg("could not determine secret type")
				}
			} else {
				// KV1
				for k, v := range s.Data {
					result[k] = v
				}
			}
		}
	}

	if r.config.Source.Debug {
		var s []string
		for k := range result {
			s = append(s, k)
		}
		r.logger.Debug().Strs("secret_keys", s).
			Msg("secret(s) found, value(s) not shown")
	}

	r.secrets = result

	return nil
}

// sanitize - sanitizes keys converting dashes(-) and dots(.) to underscores
func (r *Resource) sanitize() {
	if !r.config.Source.Sanitize {
		return
	}

	s := make(map[string]interface{}, 0)
	for k, v := range r.secrets {
		k = strings.Replace(k, "-", "_", -1)
		k = strings.Replace(k, ".", "_", -1)
		s[k] = v
	}
	r.secrets = s
}

// upcase - converts keys to UPPERCASE
func (r *Resource) upcase() {
	if !r.config.Source.Upcase {
		return
	}

	s := make(map[string]interface{}, 0)
	for k, v := range r.secrets {
		s[strings.ToUpper(k)] = v
	}
	r.secrets = s
}

// write - writes the secrets to a file
func (r Resource) write(b []byte) error {
	f, err := os.OpenFile(
		fmt.Sprintf("%s/secrets", r.workDir),
		os.O_CREATE|os.O_WRONLY, 0644,
	)
	if err != nil {
		r.logger.Fatal().Err(err).
			Msg("error opening file for write")
	}
	defer f.Close()

	if _, err := f.Write(b); err != nil {
		r.logger.Fatal().Err(err).
			Msg("error writing to destination file")
	}

	return nil
}
