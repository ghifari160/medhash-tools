package cmd_test

import (
	"path/filepath"
	"testing"

	"github.com/ghifari160/medhash-tools/cmd"
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

		testcommon.Case("default", "default"),
		testcommon.Case("all", "all"),
	}

	testcommon.RunCases(t, testGen, cases)
}

func testGen(t *testing.T, alg string, opts ...testcommon.Options) {
	t.Parallel()

	require := require.New(t)
	dir := t.TempDir()
	payload := testcommon.GenPayload(t, dir, testcommon.PayloadSize())

	c := new(cmd.Gen)
	c.Dirs = []string{dir}
	var conf medhash.Config

	switch alg {
	case "xxh3":
		c.XXH3 = true
		conf.XXH3 = true
	case "sha512":
		c.SHA512 = true
		conf.SHA512 = true
	case "sha3":
		c.SHA3 = true
		conf.SHA3 = true
	case "sha256":
		c.SHA256 = true
		conf.SHA256 = true
	case "sha1":
		c.SHA1 = true
		conf.SHA1 = true
	case "md5":
		c.MD5 = true
		conf.MD5 = true
	case "all":
		c.All = true
		conf = medhash.AllConfig
	default:
		c.Default = true
		conf = medhash.DefaultConfig
	}
	conf.Dir = dir
	conf.Manifest = medhash.DefaultManifestName

	require.Zero(c.Execute())
	require.FileExists(filepath.Join(dir, conf.Manifest))
	testcommon.VerifyManifest(t, conf, payload.Hash)
}
