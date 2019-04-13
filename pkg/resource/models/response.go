package models

// Response - the response
type Response struct {
	Metadata Metadata `json:"metadata"`
	Version  Version  `json:"version"`
}
