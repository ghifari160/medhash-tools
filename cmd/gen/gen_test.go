package gen_test

import (
	"path/filepath"
	"testing"

	"github.com/ghifari160/medhash-tools/cmd/gen"
	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/ghifari160/medhash-tools/testcommon"
	"github.com/stretchr/testify/require"
)

func TestGen(t *testing.T) {
	t.Parallel()

	cases := []testcommon.TestCase{
		testcommon.Case("xxh3", "xxh3"),
		testcommon.Case("sha512", "sha512"),
		testcommon.Case("sha3", "sha3"),
		testcommon.Case("sha256", "sha256"),
		testcommon.Case("sha1", "sha1"),
		testcommon.Case("md5", "md5"),

		testcommon.Case("default/default", "default"),
		testcommon.Case("default/signed", "default", withEd25519()),
		testcommon.Case("all", "all"),
	}

	testcommon.RunCases(t, testGen, cases)
}

func testGen(t *testing.T, alg string, opts ...testcommon.Options) {
	t.Parallel()

	require := require.New(t)
	dir := t.TempDir()
	payload := testcommon.GenPayload(t, dir, testcommon.PayloadSize())

	options := testcommon.MergeOptions(opts...)

	command := gen.CommandGen()
	var conf medhash.Config

	arguments := make([]string, 2)
	arguments[0] = "gen"

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
		conf.Ed25519.PubKey = testcommon.LoadKey(t, pubPath)
		arguments = append(arguments, "--ed25519", privPath)
	}

	conf.Dir = dir
	conf.Manifest = medhash.DefaultManifestName
	arguments = append(arguments, dir)

	err := command.Run(t.Context(), arguments)
	require.NoError(err)
	require.FileExists(filepath.Join(dir, conf.Manifest))
	testcommon.VerifyManifest(t, conf, payload.Hash)
}

func withEd25519() testcommon.Options {
	return testcommon.NewOptions("ed25519", true)
}
