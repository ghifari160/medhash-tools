package medhash

import "path/filepath"

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

	if config.SHA3 {
		if med.Hash.SHA3_256 != m.Hash.SHA3_256 {
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
