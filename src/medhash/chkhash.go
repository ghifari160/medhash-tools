// MedHash Tools
// Copyright (c) 2021 GHIFARI160
// MIT License

package medhash

func ChkHash(path string, hash *Hash) (bool, error) {
	h, err := GenHash(path)
	if err != nil {
		return false, fmtError(err)
	}

	if h.SHA256 != hash.SHA256 || h.SHA3_256 != hash.SHA3_256 || h.SHA1 != hash.SHA1 || h.MD5 != hash.MD5 {
		return false, nil
	}

	return true, nil
}
