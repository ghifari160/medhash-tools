package cmd_test

import (
	"path/filepath"
	"testing"

	"github.com/ghifari160/medhash-tools/cmd"
	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/ghifari160/medhash-tools/testcommon"
	"github.com/stretchr/testify/require"
)

func TestUpgrade(t *testing.T) {
	cases := []testcommon.TestCase{
		testcommon.Case("0.1.0", "0.1.0", withGenConfig(medhash.Config{SHA256: true})),
		testcommon.Case("0.2.0", "0.2.0", withGenConfig(medhash.Config{SHA256: true})),
		testcommon.Case("0.3.0", "0.3.0", withGenConfig(medhash.Config{
			SHA256: true,
			SHA1:   true,
			MD5:    true,
		})),
		testcommon.Case("0.4.0", "0.4.0", withGenConfig(medhash.Config{
			SHA3:   true,
			SHA256: true,
			SHA1:   true,
			MD5:    true,
		})),
		testcommon.Case("0.5.0/not_forced", "0.5.0", withGenConfig(medhash.Config{XXH3: true}), withForce(false)),
		testcommon.Case("0.5.0/forced", "0.5.0", withGenConfig(medhash.Config{XXH3: true}), withForce(true)),
	}

	testcommon.RunCases(t, testUpgrade, cases)
}

func testUpgrade(t *testing.T, version string, opts ...testcommon.Options) {
	t.Parallel()

	require := require.New(t)
	dir := t.TempDir()
	payload := testcommon.GenPayload(t, dir, testcommon.PayloadSize())
	manifestPath := filepath.Join(dir, medhash.DefaultManifestName)

	options := testcommon.MergeOptions(opts...)
	genConf := genConfig(options)
	force := options.Bool("force")

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
