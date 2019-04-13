package models

// Request - the request
type Request struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}
