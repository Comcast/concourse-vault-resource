package models

// Source - source configuration for the resource
type Source struct {
	// VaultPaths - the path(s) to the secrets in vault.
	VaultPaths map[string]int `json:"vault_paths"`

	// Format - the desired output format. Supported formats are yaml or json.
	Format string `json:"format"`

	// Prefix - a desired prefix to prepend to a secret key.
	Prefix string `json:"prefix"`

	// RoleID - the role_id for approle authentication.
	RoleID string `json:"role_id"`

	// RoleName - the role_name for approle authentication.
	RoleName string `json:"role_name"`

	// SecretID - the secret_id for approle authentication.
	SecretID string `json:"secret_id"`

	// VaultAddr - the address to the vault server.
	VaultAddr string `json:"vault_addr"`

	// VaultToken - the token to use to authenticate to vault.
	VaultToken string `json:"vault_token"`

	// Retries - the amount of times to try to read a secret from vault.
	Retries int `json:"retries"`

	// Debug - enable debug logging.
	Debug bool `json:"debug"`

	// Sanitize - convert dashes and dots to underscores in vault keys.
	Sanitize bool `json:"sanitize"`

	// Upcase - conver the vault keys to uppercase.
	Upcase bool `json:"upcase"`

	// VaultInsecure - connect the the vault server with insecure.
	VaultInsecure bool `json:"vault_insecure"`
}
