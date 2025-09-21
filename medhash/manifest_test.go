package medhash_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/ghifari160/medhash-tools/testcommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSON(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	assert := assert.New(t)
	dir := t.TempDir()
	payload := testcommon.GenPayload(t, dir, testcommon.PayloadSize())
	require.FileExists(filepath.Join(dir, payload.Path))

	config := medhash.Config{
		Dir: dir,
	}

	manifest, err := medhash.NewWithConfig(config)
	require.NoError(err)

	require.NoError(manifest.Add(payload.Path))

	expect, err := json.MarshalIndent(manifest, "", "  ")
	require.NoError(err)

	actual, err := manifest.JSON()
	require.NoError(err)

	assert.Equal(expect, actual)
}

func TestJSONStream(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	assert := assert.New(t)
	dir := t.TempDir()
	payload := testcommon.GenPayload(t, dir, testcommon.PayloadSize())
	require.FileExists(filepath.Join(dir, payload.Path))

	config := medhash.Config{
		Dir: dir,
	}

	manifestPath := filepath.Join(dir, medhash.DefaultManifestName)

	manifest, err := medhash.NewWithConfig(config)
	require.NoError(err)

	require.NoError(manifest.Add(payload.Path))

	expected, err := json.MarshalIndent(manifest, "", "  ")
	expected = append(expected, []byte("\n")...)
	require.NoError(err)

	f, err := os.OpenFile(manifestPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	require.NoError(err)
	require.NoError(manifest.JSONStream(f))
	require.NoError(f.Close())
	require.FileExists(manifestPath)

	actual, err := os.ReadFile(manifestPath)
	require.NoError(err)

	assert.Equal(expected, actual)
}

func withInvalidKey(private bool) testcommon.Options {
	if private {
		return testcommon.NewOptions("invalid_private_key", true)
	} else {
		return testcommon.NewOptions("invalid_public_key", true)
	}
}

func withMalformedSignature() testcommon.Options {
	return testcommon.NewOptions("malformed_signature", true)
}

func withBadSignature() testcommon.Options {
	return testcommon.NewOptions("bad_signature", true)
}
