package upgrade_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ghifari160/medhash-tools/cmd/upgrade"
	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/ghifari160/medhash-tools/testcommon"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
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
	chkConf := medhash.DefaultConfig
	chkConf.Dir = dir
	chkConf.Manifest = filepath.Base(manifestPath)

	var shouldError bool

	command := upgrade.CommandUpgrade()
	command.ExitErrHandler = func(ctx context.Context, c *cli.Command, err error) {
		if shouldError {
			require.Error(err)
		} else {
			require.NoError(err)
		}
	}
	arguments := make([]string, 1)
	arguments[0] = "upgrade"

	if force {
		arguments = append(arguments, "--force")
	}

	if version == "0.1.0" {
		testcommon.CreateLegacyManifest(t, dir, payload)
	} else {
		testcommon.CreateManifest(t, genConf, payload, version)
	}

	if version == upgrade.CurrentSpec && !force {
		if !force {
			shouldError = true
		}
	}

	arguments = append(arguments, dir)

	err := command.Run(t.Context(), arguments)

	if shouldError {
		require.Error(err)
	} else {
		require.NoError(err)
	}
	require.FileExists(manifestPath)
	testcommon.VerifyManifest(t, chkConf, payload.Hash)
}

func withForce(force bool) testcommon.Options {
	return testcommon.NewOptions("force", force)
}

func withGenConfig(config medhash.Config) testcommon.Options {
	return testcommon.NewOptions("config_gen", config)
}

func genConfig(options testcommon.Options) medhash.Config {
	if v, ok := options.Raw("config_gen").(medhash.Config); ok {
		return v
	} else {
		return medhash.Config{}
	}
}
