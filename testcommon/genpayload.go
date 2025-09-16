package testcommon

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/stretchr/testify/require"
	"github.com/zeebo/xxh3"
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

	payload.Hash.XXH3 = hashToString(t, xxh3)
	payload.Hash.SHA512 = hashToString(t, sha512)
	payload.Hash.SHA3 = hashToString(t, sha3)
	payload.Hash.SHA256 = hashToString(t, sha256)
	payload.Hash.SHA1 = hashToString(t, sha1)
	payload.Hash.MD5 = hashToString(t, md5)

	t.Log("Done generating payload")

	return
}

func hashToString(t testing.TB, hash hash.Hash) string {
	t.Helper()
	return hex.EncodeToString(hash.Sum(nil))
}

// PayloadSize returns the default payload size.
// If testing.Short, the payload size is 1024.
// Otherwise, it is MaxBuffer.
func PayloadSize() int64 {
	if testing.Short() {
		return 1024
	} else {
		return MaxBuffer
	}
}
