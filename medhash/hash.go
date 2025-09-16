package medhash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/zeebo/xxh3"
)

// genHash generates a hash for the media specified in the config path.
func genHash(config Config, media string) (med Media, err error) {
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

	f, err := os.Open(filepath.Join(config.Dir, media))
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

	med.Path = filepath.ToSlash(media)
	med.Hash = hash

	return
}

// chkHash verifies the hash for the media.
// Hashes for the media are verified at the same time.
// It is up to the caller to determine which hash are verified by specifying the appropriate flags
// in config.
func chkHash(config Config, med Media) (err error) {
	mediaPath := filepath.FromSlash(med.Path)

	if med.Hash.XXH3 == "" {
		config.XXH3 = false
	}
	if med.Hash.SHA512 == "" {
		config.SHA512 = false
	}
	if med.Hash.SHA3 == "" && med.Hash.SHA3_256 == "" {
		config.SHA3 = false
	}
	if med.Hash.SHA256 == "" {
		config.SHA256 = false
	}
	if med.Hash.SHA1 == "" {
		config.SHA1 = false
	}
	if med.Hash.MD5 == "" {
		config.MD5 = false
	}

	chk, err := genHash(config, mediaPath)
	if err != nil {
		return
	}

	if config.XXH3 && !hashEq(med.Hash.XXH3, chk.Hash.XXH3) {
		err = hashErr{"XXH3", med.Hash.XXH3, chk.Hash.XXH3}
	}
	if config.SHA512 && !hashEq(med.Hash.SHA512, chk.Hash.SHA512) {
		err = hashErr{"SHA512", med.Hash.SHA512, chk.Hash.SHA512}
	}
	if config.SHA3 {
		a := med.Hash.SHA3
		if a == "" {
			a = med.Hash.SHA3_256
		}
		b := chk.Hash.SHA3
		if b == "" {
			b = med.Hash.SHA3_256
		}

		if !hashEq(a, b) {
			err = hashErr{"SHA3", a, b}
		}
	}
	if config.SHA256 && !hashEq(med.Hash.SHA256, chk.Hash.SHA256) {
		err = hashErr{"SHA256", med.Hash.SHA256, chk.Hash.SHA256}
	}
	if config.SHA1 && !hashEq(med.Hash.SHA1, chk.Hash.SHA1) {
		err = hashErr{"SHA1", med.Hash.SHA1, chk.Hash.SHA1}
	}
	if config.MD5 && !hashEq(med.Hash.MD5, chk.Hash.MD5) {
		err = hashErr{"MD5", med.Hash.MD5, chk.Hash.MD5}
	}

	return
}

func hashEq(a, b string) bool {
	return a == b
}

type hashErr struct {
	alg      string
	expected string
	actual   string
}

func (err hashErr) Error() string {
	return "expected " + strings.ToUpper(err.alg) + " hash: " +
		strconv.Quote(err.expected) + " actual: " + strconv.Quote(err.actual)
}
