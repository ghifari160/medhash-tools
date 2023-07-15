package testcommon

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/stretchr/objx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CreateLegacyManifest creates a v0.1.0 Manifest for use in tests.
func CreateLegacyManifest(t testing.TB, dir string, payload medhash.Media) {
	t.Helper()

	require := require.New(t)

	f, err := os.Create(filepath.Join(dir, "sums.txt"))
	require.NoError(err)

	_, err = f.WriteString(fmt.Sprintf("%s        %s\n", payload.Path, payload.Hash.SHA256))
	require.NoError(err)

	err = f.Close()
	require.NoError(err)
}

// CreateManifest creates a Manifest version ver for use in tests.
// It is up to the caller to follow the spec of the specified version.
func CreateManifest(t testing.TB, dir string, payload medhash.Media, ver string,
	config medhash.Config) {
	t.Helper()

	require := require.New(t)

	if !config.SHA3 {
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

	man := medhash.New()
	man.Version = ver
	man.Generator = "MedHash Tools Test"

	man.Media = []medhash.Media{payload}

	manFile, err := json.Marshal(man)
	require.NoError(err)

	f, err := os.Create(filepath.Join(dir, medhash.DefaultManifestName))
	require.NoError(err)

	_, err = f.Write(manFile)
	require.NoError(err)

	err = f.Close()
	require.NoError(err)
}

// VerifyManifest verifies a Manifest in dir.
// It will only verify the first Media.
// The hashes of the Media are verified against hash.
func VerifyManifest(t testing.TB, dir string, config medhash.Config, hash medhash.Hash) {
	t.Helper()

	require := require.New(t)
	assert := assert.New(t)

	manFile, err := os.ReadFile(filepath.Join(dir, medhash.DefaultManifestName))
	require.NoError(err)

	manifest, err := objx.FromJSON(string(manFile))
	require.NoError(err)

	require.Equal(medhash.ManifestFormatVer, manifest.Get("version").Str())

	if config.SHA3 {
		assert.Equal(hash.SHA3_256, manifest.Get("media[0].hash.sha3-256").Str())
	}

	if config.SHA256 {
		assert.Equal(hash.SHA256, manifest.Get("media[0].hash.sha256").Str())
	}

	if config.SHA1 {
		assert.Equal(hash.SHA1, manifest.Get("media[0].hash.sha1").Str())
	}

	if config.MD5 {
		assert.Equal(hash.MD5, manifest.Get("media[0].hash.md5").Str())
	}
}
