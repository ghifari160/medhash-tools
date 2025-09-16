package cmd_test

import (
	"testing"

	"github.com/ghifari160/medhash-tools/cmd"
	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/ghifari160/medhash-tools/testcommon"
	"github.com/stretchr/testify/require"
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

	c := new(cmd.Chk)
	c.Dirs = []string{dir}
	if len(files) > 0 {
		c.Files = files
	}
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

	if invalidate {
		payload.Hash.XXH3 = "__INVALID__"
		payload.Hash.SHA512 = "__INVALID__"
		payload.Hash.SHA3 = "__INVALID__"
		payload.Hash.SHA3_256 = "__INVALID__"
		payload.Hash.SHA256 = "__INVALID__"
		payload.Hash.SHA1 = "__INVALID__"
		payload.Hash.MD5 = "__INVALID__"
	}

	testcommon.CreateManifest(t, conf, payload, medhash.ManifestFormatVer)

	if !invalidate {
		require.Zero(c.Execute())
	} else {
		require.NotZero(c.Execute())
	}
}
