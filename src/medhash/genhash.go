// MedHash Tools
// Copyright (c) 2021 GHIFARI160
// MIT License

package medhash

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/sha3"
)

const DEFAULT_BUFFERSIZE int = 4096

// Deprecated: legacy code.
func bufferedGenHash(path string, hasher *hash.Hash, bufferSize int) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	file := bufio.NewReaderSize(f, bufferSize)

	_, err = io.Copy((*hasher), file)
	if err != nil {
		return err
	}

	return f.Close()
}

// Deprecated: legacy code.
func GenHash(path string) (*Media, error) {
	var err error

	var hasher hash.Hash
	hash := Hash{}

	hasher = sha256.New()
	err = bufferedGenHash(path, &hasher, DEFAULT_BUFFERSIZE)
	if err != nil {
		return nil, fmtError(err)
	}
	hash.SHA256 = hex.EncodeToString(hasher.Sum(nil))

	hasher = sha3.New256()
	err = bufferedGenHash(path, &hasher, DEFAULT_BUFFERSIZE)
	if err != nil {
		return nil, fmtError(err)
	}
	hash.SHA3_256 = hex.EncodeToString(hasher.Sum(nil))

	hasher = sha1.New()
	err = bufferedGenHash(path, &hasher, DEFAULT_BUFFERSIZE)
	if err != nil {
		return nil, fmtError(err)
	}
	hash.SHA1 = hex.EncodeToString(hasher.Sum(nil))

	hasher = md5.New()
	err = bufferedGenHash(path, &hasher, DEFAULT_BUFFERSIZE)
	if err != nil {
		return nil, fmtError(err)
	}
	hash.MD5 = hex.EncodeToString(hasher.Sum(nil))

	return &Media{
		Path: filepath.ToSlash(path),
		Hash: &hash,
	}, nil
}
