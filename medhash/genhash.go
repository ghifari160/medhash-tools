package medhash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"

	"github.com/zeebo/xxh3"
	"golang.org/x/crypto/sha3"
)

// GenHash generates a hash for the media specified in the config path.
// Hashes for the media are generated at the same time.
func GenHash(config Config) (med Media, err error) {
	return genHash(config)
}

// genHash generates a hash for the media specified in the config path.
func genHash(config Config) (med Media, err error) {
	writers := make([]io.Writer, 0)
	hashers := make(map[string]hash.Hash)

	if config.XXH3 {
		hashers["xxh3"] = xxh3.New()
		writers = append(writers, hashers["xxh3"])
	}

	if config.SHA512 {
		hashers["sha512"] = sha512.New()
		writers = append(writers, hashers["sha512"])
	}

	if config.SHA3 {
		hashers["sha3"] = sha3.New256()
		writers = append(writers, hashers["sha3"])
	}

	if config.SHA256 {
		hashers["sha256"] = sha256.New()
		writers = append(writers, hashers["sha256"])
	}

	if config.SHA1 {
		hashers["sha1"] = sha1.New()
		writers = append(writers, hashers["sha1"])
	}

	if config.MD5 {
		hashers["md5"] = md5.New()
		writers = append(writers, hashers["md5"])
	}

	writer := io.MultiWriter(writers...)

	f, err := os.Open(filepath.Join(config.Dir, config.Path))
	if err != nil {
		return
	}
	defer f.Close()

	_, err = io.Copy(writer, f)
	if err != nil {
		return
	}

	hash := Hash{}

	if h := hashers["xxh3"]; h != nil {
		hash.XXH3 = hex.EncodeToString(h.Sum(nil))
	}

	if h := hashers["sha512"]; h != nil {
		hash.SHA512 = hex.EncodeToString(h.Sum(nil))
	}

	if h := hashers["sha3"]; h != nil {
		hash.SHA3 = hex.EncodeToString(h.Sum(nil))
		hash.SHA3_256 = hex.EncodeToString(h.Sum(nil))
	}

	if h := hashers["sha256"]; h != nil {
		hash.SHA256 = hex.EncodeToString(h.Sum(nil))
	}

	if h := hashers["sha1"]; h != nil {
		hash.SHA1 = hex.EncodeToString(h.Sum(nil))
	}

	if h := hashers["md5"]; h != nil {
		hash.MD5 = hex.EncodeToString(h.Sum(nil))
	}

	med.Path = filepath.ToSlash(config.Path)
	med.Hash = hash

	return
}
