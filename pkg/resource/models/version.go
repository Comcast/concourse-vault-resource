package models

// Version - the version of the resource
type Version struct {
	// Path - the path to the secret in vault
	Path string `json:"path"`
	// Version - the version of the resource.
	Version string `json:"version"`
}
