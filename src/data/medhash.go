package data

// Medhash object definition
//
// Deprecated: legacy code.
type Medhash struct {
	Version string  `json:"version"`
	Media   []Media `json:"media"`
}
