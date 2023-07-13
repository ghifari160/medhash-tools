// MedHash Tools
// Copyright (c) 2023 GHIFARI160
// MIT License

package medhash

const ManifestFormatVer = "0.4.0"
const DefaultManifestName = "medhash.json"

// Manifest is a MedHash manifest.
type Manifest struct {
	Version   string  `json:"version"`
	Generator string  `json:"generator,omitempty"`
	Media     []Media `json:"media"`

	Config Config `json:"-"`
}

// New creates a new Manifest with the default configuration.
func New() *Manifest {
	return NewWithConfig(Config{
		SHA3:   true,
		SHA256: true,
		SHA1:   true,
		MD5:    true,
	})
}

// NewWithConfig creates a new Manifest with the specific configuration.
func NewWithConfig(config Config) *Manifest {
	return &Manifest{
		Version: ManifestFormatVer,
		Config:  config,
	}
}

// Media stores metadata about the media.
type Media struct {
	Path string `json:"path"`
	Hash Hash   `json:"hash"`
}

// Hash stores each hash of a Media.
type Hash struct {
	SHA256   string `json:"sha256,omitempty"`
	SHA3_256 string `json:"sha3-256,omitempty"`

	// Deprecated: SHA1 support is deprecated in spec v0.4.0.
	SHA1 string `json:"sha1,omitempty"`
	// Deprecated: MD5 support is deprecated in spec v0.4.0.
	MD5 string `json:"md5,omitempty"`
}
