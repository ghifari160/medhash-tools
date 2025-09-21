package chk_test

import (
	"context"
	"testing"

	"github.com/ghifari160/medhash-tools/cmd/chk"
	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/ghifari160/medhash-tools/testcommon"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestChk(t *testing.T) {
	t.Parallel()

	cases := []testcommon.TestCase{
		testcommon.Case("xxh3", "xxh3"),
		testcommon.Case("sha512", "sha512"),
		testcommon.Case("sha3", "sha3"),
		testcommon.Case("sha256", "sha256"),
		testcommon.Case("sha1", "sha1"),
		testcommon.Case("md5", "md5"),

		testcommon.Case("all", "all"),
		testcommon.Case("default/default", "default"),
		testcommon.Case("default/signed", "default", withEd25519()),
		testcommon.Case("default/invalid", "default", withInvalidate(true)),
		testcommon.Case("default/file_list/skip", "default", withFiles([]string{"payload2"})),
		testcommon.Case("default/file_list/include", "default", withFiles([]string{"payload"})),
	}

	testcommon.RunCases(t, testChk, cases)
}

func testChk(t *testing.T, alg string, opts ...testcommon.Options) {
	t.Parallel()

	require := require.New(t)
	dir := t.TempDir()
	payload := testcommon.GenPayload(t, dir, testcommon.PayloadSize())

	options := testcommon.MergeOptions(opts...)
	invalidate := options.Bool("invalidate")
	files := options.StrSlice("files")

	var shouldError bool

	command := chk.CommandChk()
	command.ExitErrHandler = func(ctx context.Context, c *cli.Command, err error) {
		if shouldError {
			require.Error(err)
		} else {
			require.NoError(err)
		}
	}
	var conf medhash.Config

	arguments := make([]string, 2)
	arguments[0] = "chk"

	if len(files) > 0 {
		for _, file := range files {
			arguments = append(arguments, "--file", file)
		}
	}

	switch alg {
	case "xxh3":
		conf.XXH3 = true
		arguments[1] = "--xxh3"
	case "sha512":
		conf.SHA512 = true
		arguments[1] = "--sha512"
	case "sha3":
		conf.SHA3 = true
		arguments[1] = "--sha3"
	case "sha256":
		conf.SHA256 = true
		arguments[1] = "--sha256"
	case "sha1":
		conf.SHA1 = true
		arguments[1] = "--sha1"
	case "md5":
		conf.MD5 = true
		arguments[1] = "--md5"
	case "all":
		conf = medhash.AllConfig
		arguments[1] = "--all"
	default:
		conf = medhash.DefaultConfig
		arguments[1] = "--default"
	}

	if options.Bool("ed25519") {
		pubPath, privPath := testcommon.GenEd25519Keypair(t)
		conf.Ed25519.Enabled = true
		conf.Ed25519.PrivKey = testcommon.LoadKey(t, privPath)
		conf.Ed25519.PubKey = testcommon.LoadKey(t, pubPath)
		arguments = append(arguments, "--ed25519", pubPath)
	}

	conf.Dir = dir
	conf.Manifest = medhash.DefaultManifestName
	arguments = append(arguments, dir)

	if invalidate {
		payload.Hash.XXH3 = "__INVALID__"
		payload.Hash.SHA512 = "__INVALID__"
		payload.Hash.SHA3 = "__INVALID__"
		payload.Hash.SHA3_256 = "__INVALID__"
		payload.Hash.SHA256 = "__INVALID__"
		payload.Hash.SHA1 = "__INVALID__"
		payload.Hash.MD5 = "__INVALID__"
		shouldError = true
	}

	testcommon.CreateManifest(t, conf, payload, medhash.ManifestFormatVer)

	err := command.Run(t.Context(), arguments)
	if !invalidate {
		require.NoError(err)
	} else {
		require.Error(err)
	}
}

// withInvalidate invalidates the payload hash for testing.
func withInvalidate(invalidate bool) testcommon.Options {
	return testcommon.NewOptions("invalidate", invalidate)
}

// withFiles specifies the Files argument to a command for testing.
func withFiles(files []string) testcommon.Options {
	return testcommon.NewOptions("files", files)
}

// withEd25519 flags the test to generate and verify a signed Manifest.
func withEd25519() testcommon.Options {
	return testcommon.NewOptions("ed25519", true)
}
