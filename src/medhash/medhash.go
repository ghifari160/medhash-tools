// MedHash Tools
// Copyright (c) 2021 GHIFARI160
// MIT License

package medhash

import "errors"

const MEDHASH_FORMAT_VERSION_PREVIOUS = "0.3.0"
const MEDHASH_FORMAT_VERSION_CURRENT = "0.4.0"
const MEDHASH_MANIFEST_NAME = "medhash.json"
const MEDHASH_ERROR_PREFIX = "medhash:"

type MedHash struct {
	Version   string  `json:"version"`
	Generator string  `json:"generator,omitempty"`
	Media     []Media `json:"media"`
}

func New() *MedHash {
	return &MedHash{Version: MEDHASH_FORMAT_VERSION_CURRENT}
}

type Media struct {
	Path string `json:"path"`
	Hash *Hash  `json:"hash"`
}

type Hash struct {
	SHA256   string `json:"sha256"`
	SHA3_256 string `json:"sha3-256"`
	SHA1     string `json:"sha1"`
	MD5      string `json:"md5"`
}

func fmtError(err error) error {
	if err != nil {
		errMsg := MEDHASH_ERROR_PREFIX + " " + err.Error()

		return errors.New(errMsg)
	}

	return nil
}
