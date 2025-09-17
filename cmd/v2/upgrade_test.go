package cmd_test

import (
	"path/filepath"
	"testing"

	"github.com/ghifari160/medhash-tools/cmd/v2"
	"github.com/ghifari160/medhash-tools/medhash/v2"
	"github.com/ghifari160/medhash-tools/testcommon/v2"
	"github.com/stretchr/testify/require"
)

func TestUpgrade(t *testing.T) {
	type testCase struct {
		id      string
		version string
		genConf medhash.Config
		force   bool
	}

	cases := []testCase{
		{"0.1.0", "0.1.0", medhash.Config{SHA256: true}, false},
		{"0.2.0", "0.2.0", medhash.Config{SHA256: true}, false},
		{"0.3.0", "0.3.0", medhash.Config{
			SHA256: true,
			SHA1:   true,
			MD5:    true,
		}, false},
		{"0.4.0", "0.4.0", medhash.Config{
			SHA3:   true,
			SHA256: true,
			SHA1:   true,
			MD5:    true,
		}, false},
		{"0.5.0/not_forced", "0.5.0", medhash.Config{XXH3: true}, false},
		{"0.5.0/forced", "0.5.0", medhash.Config{XXH3: true}, true},
	}

	for _, testCase := range cases {
		t.Run(testCase.id, func(t *testing.T) {
			testUpgrade(t, testCase.genConf, testCase.version, testCase.force)
		})
	}
}

func testUpgrade(t *testing.T, genConf medhash.Config, version string, force bool) {
	t.Parallel()

	require := require.New(t)
	dir := t.TempDir()
	payload := testcommon.GenPayload(t, dir, testcommon.PayloadSize())
	manifestPath := filepath.Join(dir, medhash.DefaultManifestName)

	genConf.Dir = dir
	if genConf.Manifest == "" {
		genConf.Manifest = medhash.DefaultManifestName
	}
	checkConf := medhash.DefaultConfig
	checkConf.Dir = dir
	checkConf.Manifest = filepath.Base(manifestPath)

	c := new(cmd.Upgrade)
	c.Default = true
	c.Dirs = []string{dir}

	if version == "0.1.0" {
		testcommon.CreateLegacyManifest(t, dir, payload)
	} else {
		testcommon.CreateManifest(t, genConf, payload, version)
	}

	if version == cmd.CurrentSpec {
		if !force {
			require.NotZero(c.Execute())
			return
		} else {
			c.Force = true
		}
	}

	require.Zero(c.Execute())
	require.FileExists(manifestPath)
	testcommon.VerifyManifest(t, checkConf, payload.Hash)
}
