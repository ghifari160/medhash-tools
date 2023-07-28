package medhash

const ManifestFormatVer = "0.5.0"
const DefaultManifestName = "medhash.json"

var (
	DefaultConfig = Config{
		XXH3: true,
	}
	AllConfig = Config{
		XXH3:   true,
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
	Version   string  `json:"version"`
	Generator string  `json:"generator,omitempty"`
	Media     []Media `json:"media"`

	Config Config `json:"-"`
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

// Hash stores each hash of a Media.
type Hash struct {
	XXH3     string `json:"xxh3,omitempty"`
	SHA256   string `json:"sha256,omitempty"`
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
	// SHA3 toggles the SHA3-256 hash generation.
	SHA3 bool
	// SHA256 toggles the SHA256 hash generation.
	SHA256 bool
	// SHA1 toggles the SHA1 hash generation.
	SHA1 bool
	// MD5 toggles the MD5 hash generation.
	MD5 bool
}
