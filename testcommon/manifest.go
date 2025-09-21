package testcommon

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CreateLegacyManifest creates a v0.1.0 Manifest for use in tests.
func CreateLegacyManifest(t testing.TB, dir string, payload medhash.Media) {
	t.Helper()

	require := require.New(t)

	f, err := os.Create(filepath.Join(dir, "sums.txt"))
	require.NoError(err)

	_, err = fmt.Fprintf(f, "%s        %s\n", payload.Path, payload.Hash.SHA256)
	require.NoError(err)

	err = f.Close()
	require.NoError(err)
}

// CreateManifest creates a Manifest version ver for use in tests.
// It is up to the caller to follow the spec of the specified version.
func CreateManifest(t testing.TB, config medhash.Config, payload medhash.Media, ver string) {
	t.Helper()
	require := require.New(t)
	manifestPath := filepath.Join(config.Dir, config.Manifest)

	if !config.XXH3 {
		payload.Hash.XXH3 = ""
	}
	if !config.SHA512 {
		payload.Hash.SHA512 = ""
	}
	if !config.SHA3 {
		payload.Hash.SHA3 = ""
		payload.Hash.SHA3_256 = ""
	}
	if !config.SHA256 {
		payload.Hash.SHA256 = ""
	}
	if !config.SHA1 {
		payload.Hash.SHA1 = ""
	}
	if !config.MD5 {
		payload.Hash.MD5 = ""
	}

	manifest, err := medhash.NewWithConfig(config)
	require.NoError(err)
	manifest.Version = ver
	manifest.Generator = "MedHash Tools Test"
	manifest.Media = []medhash.Media{payload}

	if config.Ed25519.Enabled {
		if config.Ed25519.PrivKey == nil {
			t.Fatalf("Ed25519.Enabled is true but no valid PrivKey found")
		}

		payload, err := manifest.JSON()
		require.NoError(err)

		require.NotPanics(func() {
			signature := ed25519.Sign(config.Ed25519.PrivKey, payload)
			manifest.Signature.Ed25519 = hex.EncodeToString(signature)
		})
	}

	require.NoError(storeManifest(manifest, manifestPath))
	require.FileExists(manifestPath)
}

// VerifyManifest verifies a Manifest in dir.
// It will only verify the first Media.
// The hashes of the Media are verified against hash.
func VerifyManifest(t testing.TB, config medhash.Config, hash medhash.Hash) {
	t.Helper()
	require := require.New(t)
	assert := assert.New(t)
	manifestPath := filepath.Join(config.Dir, config.Manifest)

	if hash.XXH3 == "" {
		config.XXH3 = false
	}
	if hash.SHA512 == "" {
		config.SHA512 = false
	}
	if hash.SHA3 == "" && hash.SHA3_256 == "" {
		config.SHA3 = false
	}
	if hash.SHA256 == "" {
		config.SHA256 = false
	}
	if hash.SHA1 == "" {
		config.SHA1 = false
	}
	if hash.MD5 == "" {
		config.MD5 = false
	}

	require.FileExists(manifestPath)
	manifest, err := loadManifest(manifestPath)
	require.NoError(err)

	for _, media := range manifest.Media {
		if config.XXH3 {
			assert.Equal(hash.XXH3, media.Hash.XXH3)
		}
		if config.SHA512 {
			assert.Equal(hash.SHA512, media.Hash.SHA512)
		}
		if config.SHA3 {
			var expected, actual string
			if hash.SHA3 != "" {
				expected = hash.SHA3
			} else {
				expected = hash.SHA3_256
			}
			if media.Hash.SHA3 != "" {
				actual = hash.SHA3
			} else {
				actual = hash.SHA3_256
			}
			assert.Equal(expected, actual)
		}
		if config.SHA256 {
			assert.Equal(hash.SHA256, media.Hash.SHA256)
		}
		if config.SHA1 {
			assert.Equal(hash.SHA1, media.Hash.SHA1)
		}
		if config.MD5 {
			assert.Equal(hash.MD5, media.Hash.MD5)
		}
	}
}

func loadManifest(path string) (manifest *medhash.Manifest, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	manifest = new(medhash.Manifest)
	err = json.NewDecoder(f).Decode(manifest)
	return
}

func storeManifest(manifest *medhash.Manifest, path string) (err error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	err = manifest.JSONStream(f)

	return
}
