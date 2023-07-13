// MedHash Tools
// Copyright (c) 2023 GHIFARI160
// MIT License

package medhash

import "path/filepath"

// ChkHash verifies the hash for the media.
// Hashes for the media are verified at the same time.
func ChkHash(dir string, med Media) (valid bool, err error) {
	var config Config

	config.Dir = dir
	config.Path = filepath.FromSlash(med.Path)

	if med.Hash.SHA3_256 != "" {
		config.SHA3 = true
	}

	if med.Hash.SHA256 != "" {
		config.SHA256 = true
	}

	if med.Hash.SHA1 != "" {
		config.SHA1 = true
	}

	if med.Hash.MD5 != "" {
		config.MD5 = true
	}

	m, err := genHash(config)
	if err != nil {
		return
	}

	valid = med.Hash.SHA3_256 == m.Hash.SHA3_256 &&
		med.Hash.SHA256 == m.Hash.SHA256 &&
		med.Hash.SHA1 == m.Hash.SHA1 &&
		med.Hash.MD5 == m.Hash.MD5

	return
}
