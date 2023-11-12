package medhash

import (
	"encoding/json"
)

const ManifestFormatVer = "0.6.0"
const DefaultManifestName = "medhash.json"

var (
	DefaultConfig = Config{
		XXH3: true,
	}
	AllConfig = Config{
		XXH3:   true,
		SHA512: true,
		SHA3:   true,
		SHA256: true,
		SHA1:   true,
		MD5:    true,
	}
	LegacyConfig = Config{
		SHA3:   true,
		SHA256: true,
		SHA1:   true,
		MD5:    true,
	}
)

// Manifest is a MedHash manifest.
type Manifest struct {
	Version   string     `json:"version"`
	Generator string     `json:"generator,omitempty"`
	Media     []Media    `json:"media"`
	Signature *Signature `json:"signature,omitempty"`

	Config Config `json:"-"`
}

// JSON serializes m to JSON.
func (m *Manifest) JSON() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}

// String implements the Stringer interface.
func (m *Manifest) String() string {
	j, err := m.JSON()
	if err != nil {
		panic(err)
	}

	return string(j)
}

// Copy copies m to target.
func (m *Manifest) Copy(target *Manifest) {
	target.Version = m.Version
	target.Generator = m.Generator
	target.Config = m.Config

	target.Media = make([]Media, len(m.Media))
	copy(target.Media, m.Media)

	if m.Signature != nil {
		target.Signature = new(Signature)
		target.Signature.Ed25519 = m.Signature.Ed25519
	}
}

// StripSignature returns a copy of m with Signature stripped.
// It also returns the stripped Signature.
func (m *Manifest) StripSignature() (stripped *Manifest, signature *Signature) {
	stripped = new(Manifest)
	m.Copy(stripped)
	stripped.Signature = nil

	signature = m.Signature

	return
}

// New creates a new Manifest with the default configuration.
func New() *Manifest {
	return NewWithConfig(DefaultConfig)
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

// Signature stores the signature of the Manifest.
type Signature struct {
	Ed25519 string `json:"ed25519,omitempty"`
}

// Hash stores each hash of a Media.
type Hash struct {
	XXH3   string `json:"xxh3,omitempty"`
	SHA512 string `json:"sha512,omitempty"`
	SHA256 string `json:"sha256,omitempty"`
	SHA3   string `json:"sha3,omitempty"`
	// Deprecated: use SHA3.
	SHA3_256 string `json:"sha3-256,omitempty"`
	SHA1     string `json:"sha1,omitempty"`
	MD5      string `json:"md5,omitempty"`
}

// Config configures the hasher.
type Config struct {
	// Dir is the path to the target directory.
	Dir string
	// Path is the path of the current media.
	Path string

	// XXH3 toggles the XXH3_64 hash generation.
	XXH3 bool
	// SHA512 toggles the SHA512 hash generation.
	SHA512 bool
	// SHA3 toggles the SHA3-256 hash generation.
	SHA3 bool
	// SHA256 toggles the SHA256 hash generation.
	SHA256 bool
	// SHA1 toggles the SHA1 hash generation.
	SHA1 bool
	// MD5 toggles the MD5 hash generation.
	MD5 bool

	// Ed25519 configures the Ed25519 signature generation/verification.
	Ed25519 struct {
		Enabled    bool
		PrivateKey []byte
		PublicKey  []byte
	}
}
