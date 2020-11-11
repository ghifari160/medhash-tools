package data

// Medhash object definition
type Medhash struct {
	Version string  `json:"version"`
	Media   []Media `json:"media"`
}
