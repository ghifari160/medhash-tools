package data

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"os"
)

func bufferedGenHash(path string, hasher *hash.Hash) error {
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy((*hasher), file)

	if err != nil {
		return err
	}

	return nil
}

// GenHash hash generator
func GenHash(path string) (Hash, error) {
	hash := Hash{}

	hasher := sha256.New()

	err := bufferedGenHash(path, &hasher)
	if err != nil {
		return hash, errors.New("Unable to generate SHA256 hash for: " + path)
	}

	hash.SHA256 = hex.EncodeToString(hasher.Sum(nil))

	hasher = sha1.New()

	err = bufferedGenHash(path, &hasher)
	if err != nil {
		return hash, errors.New("Unable to generate SHA1 hash for: " + path)
	}

	hash.SHA1 = hex.EncodeToString(hasher.Sum(nil))

	hasher = md5.New()

	err = bufferedGenHash(path, &hasher)
	if err != nil {
		return hash, errors.New("Unable to generate MD5 hash for: " + path)
	}

	hash.MD5 = hex.EncodeToString(hasher.Sum(nil))

	return hash, nil
}
