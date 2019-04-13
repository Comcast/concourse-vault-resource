package models

// MetadataKvP - for any available information from the resource
type MetadataKvP struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Metadata - metadata kvps
type Metadata []MetadataKvP
