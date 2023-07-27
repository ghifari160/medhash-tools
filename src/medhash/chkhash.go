// MedHash Tools
// Copyright (c) 2021 GHIFARI160
// MIT License

// Deprecated: legacy code.
package medhash

import "path/filepath"

// Deprecated: legacy code.
func ChkHash(media *Media) (bool, error) {
	m, err := GenHash(filepath.FromSlash(media.Path))
	if err != nil {
		return false, fmtError(err)
	}

	if m.Hash.SHA256 != media.Hash.SHA256 ||
		(media.Hash.SHA3_256 != "" && m.Hash.SHA3_256 != media.Hash.SHA3_256) ||
		m.Hash.SHA1 != media.Hash.SHA1 ||
		m.Hash.MD5 != media.Hash.MD5 {
		return false, nil
	}

	return true, nil
}
