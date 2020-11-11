package data

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/ioutil"
)

// GenHash hash generator
func GenHash(path string) (Hash, error) {
	hash := Hash{}

	hasher := sha256.New()

	s, err := ioutil.ReadFile(path)
	hasher.Write(s)
	if err != nil {
		return hash, errors.New("Unable to generate SHA256 hash for: " + path)
	}

	hash.SHA256 = hex.EncodeToString(hasher.Sum(nil))

	hasher = sha1.New()

	s, err = ioutil.ReadFile(path)
	hasher.Write(s)
	if err != nil {
		return hash, errors.New("Unable to generate SHA1 hash for: " + path)
	}

	hash.SHA1 = hex.EncodeToString(hasher.Sum(nil))

	hasher = md5.New()

	s, err = ioutil.ReadFile(path)
	hasher.Write(s)
	if err != nil {
		return hash, errors.New("Unable to generate MD5 hash for: " + path)
	}

	hash.MD5 = hex.EncodeToString(hasher.Sum(nil))

	return hash, nil
}
