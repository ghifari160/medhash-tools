package testcommon

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/stretchr/testify/require"
	"github.com/zeebo/xxh3"
	"golang.org/x/crypto/sha3"
)

const MaxBuffer = 1 * 1024 * 1024 * 1024

// GenPayload generates a test payload of size in dir.
// The payload is wrapped around medhash.Media with all supported hash precalculated.
func GenPayload(t testing.TB, dir string, size int64) (payload medhash.Media) {
	var buf []byte
	var counter int64

	require := require.New(t)

	t.Log("Generating payload")

	if size > MaxBuffer {
		buf = make([]byte, MaxBuffer)
	} else {
		buf = make([]byte, size)
	}

	payload.Path = "payload"

	f, err := os.Create(filepath.Join(dir, payload.Path))
	require.NoError(err)

	xxh3 := xxh3.New()
	sha512 := sha512.New()
	sha3 := sha3.New256()
	sha256 := sha256.New()
	sha1 := sha1.New()
	md5 := md5.New()
	writer := io.MultiWriter(f, xxh3, sha512, sha3, sha256, sha1, md5)

	for counter < size {
		n, err := rand.Read(buf)
		require.NoError(err)

		if counter+int64(n) > size {
			n, err = writer.Write(buf[:size-counter])
		} else {
			n, err = writer.Write(buf)
		}
		require.NoError(err)

		counter += int64(n)
	}

	err = f.Close()
	require.NoError(err)

	payload.Hash.XXH3 = hex.EncodeToString(xxh3.Sum(nil))
	payload.Hash.SHA512 = hex.EncodeToString(sha512.Sum(nil))
	payload.Hash.SHA3 = hex.EncodeToString(sha3.Sum(nil))
	payload.Hash.SHA3_256 = hex.EncodeToString(sha3.Sum(nil))
	payload.Hash.SHA256 = hex.EncodeToString(sha256.Sum(nil))
	payload.Hash.SHA1 = hex.EncodeToString(sha1.Sum(nil))
	payload.Hash.MD5 = hex.EncodeToString(md5.Sum(nil))

	t.Log("Done generating payload")

	return
}
