package data

// ChkHash hash checker
func ChkHash(path string, hash Hash) bool {
	h, err := GenHash(path)

	if err != nil {
		return false
	}

	if h.SHA256 != hash.SHA256 || h.SHA1 != hash.SHA1 || h.MD5 != hash.MD5 {
		return false
	}

	return true
}
