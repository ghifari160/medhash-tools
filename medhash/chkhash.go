package medhash

import (
	"path/filepath"
)

// ChkHash verifies the hash for the media.
// Hashes for the media are verified at the same time.
// It is up to the caller to determine which hash are verified by specifying the appropriate flags
// in config.
func ChkHash(config Config, med Media) (valid bool, err error) {
	config.Path = filepath.FromSlash(med.Path)

	m, err := genHash(config)
	if err != nil {
		return
	}

	if med.Hash.XXH3 == "" {
		config.XXH3 = false
	}

	if med.Hash.SHA512 == "" {
		config.SHA512 = false
	}

	if med.Hash.SHA3 == "" && med.Hash.SHA3_256 == "" {
		config.SHA3 = false
	}

	if med.Hash.SHA256 == "" {
		config.SHA256 = false
	}

	if med.Hash.SHA1 == "" {
		config.SHA1 = false
	}

	if med.Hash.MD5 == "" {
		config.MD5 = false
	}

	if config.XXH3 {
		if med.Hash.XXH3 != m.Hash.XXH3 {
			valid = false
			return
		}
	}

	if config.SHA512 {
		if med.Hash.SHA512 != m.Hash.SHA512 {
			valid = false
			return
		}
	}

	if config.SHA3 {
		if med.Hash.SHA3 != "" && med.Hash.SHA3 != m.Hash.SHA3 {
			valid = false
			return
		} else if med.Hash.SHA3 == "" && med.Hash.SHA3_256 != m.Hash.SHA3_256 {
			valid = false
			return
		}
	}

	if config.SHA256 {
		if med.Hash.SHA256 != m.Hash.SHA256 {
			valid = false
			return
		}
	}

	if config.SHA1 {
		if med.Hash.SHA1 != m.Hash.SHA1 {
			valid = false
			return
		}
	}

	if config.MD5 {
		if med.Hash.MD5 != m.Hash.MD5 {
			valid = false
			return
		}
	}

	valid = true

	return
}
