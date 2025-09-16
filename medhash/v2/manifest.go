package medhash

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
	Version   string  `json:"version"`
	Generator string  `json:"generator,omitempty"`
	Media     []Media `json:"media"`

	Config Config `json:"-"`
}

func New() (man *Manifest, err error) {
	return NewWithConfig(DefaultConfig)
}

func NewWithConfig(config Config) (man *Manifest, err error) {
	if config.Manifest == "" {
		config.Manifest = DefaultManifestName
	}

	man = &Manifest{
		Version: ManifestFormatVer,
		Config:  config,
		Media:   make([]Media, 0),
	}

	return
}

// Config configures the hasher.
type Config struct {
	// Dir is the path to the target directory.
	Dir string
	// Manifest is the manifest file name.
	Manifest string

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
}
