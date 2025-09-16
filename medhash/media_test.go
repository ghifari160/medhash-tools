package medhash_test

import (
	"path/filepath"
	"testing"

	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/ghifari160/medhash-tools/testcommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenHash(t *testing.T) {
	t.Parallel()

	t.Run("xxh3", func(t *testing.T) {
		testGenHash(t, "xxh3", func(t testing.TB, a *assert.Assertions,
			man *medhash.Manifest, pld medhash.Media) {
			a.NotEmpty(man.Media[0].Hash.XXH3)
			a.Equal(pld.Hash.XXH3, man.Media[0].Hash.XXH3)
		})
	})

	t.Run("sha512", func(t *testing.T) {
		testGenHash(t, "sha512", func(t testing.TB, a *assert.Assertions,
			man *medhash.Manifest, pld medhash.Media) {
			a.NotEmpty(man.Media[0].Hash.SHA512)
			a.Equal(pld.Hash.SHA512, man.Media[0].Hash.SHA512)
		})
	})

	t.Run("sha3", func(t *testing.T) {
		testGenHash(t, "sha3", func(t testing.TB, a *assert.Assertions,
			man *medhash.Manifest, pld medhash.Media) {
			a.NotEmpty(man.Media[0].Hash.SHA3)
			a.Equal(pld.Hash.SHA3, man.Media[0].Hash.SHA3)
		})
	})

	t.Run("sha256", func(t *testing.T) {
		testGenHash(t, "sha256", func(t testing.TB, a *assert.Assertions,
			man *medhash.Manifest, pld medhash.Media) {
			a.NotEmpty(man.Media[0].Hash.SHA256)
			a.Equal(pld.Hash.SHA256, man.Media[0].Hash.SHA256)
		})
	})

	t.Run("sha1", func(t *testing.T) {
		testGenHash(t, "sha1", func(t testing.TB, a *assert.Assertions,
			man *medhash.Manifest, pld medhash.Media) {
			a.NotEmpty(man.Media[0].Hash.SHA1)
			a.Equal(pld.Hash.SHA1, man.Media[0].Hash.SHA1)

		})
	})

	t.Run("md5", func(t *testing.T) {
		testGenHash(t, "md5", func(t testing.TB, a *assert.Assertions,
			man *medhash.Manifest, pld medhash.Media) {
			a.NotEmpty(man.Media[0].Hash.MD5)
			a.Equal(pld.Hash.MD5, man.Media[0].Hash.MD5)
		})
	})
}

func TestCheckHash(t *testing.T) {
	t.Parallel()

	t.Run("xxh3", func(t *testing.T) {
		testCheckHashValid(t, "xxh3")
	})

	t.Run("xxh3_invalid", func(t *testing.T) {
		testCheckHashInvalid(t, "xxh3")
	})

	t.Run("xxh3_empty", func(t *testing.T) {
		testCheckHashEmpty(t, "xxh3")
	})

	t.Run("sha512", func(t *testing.T) {
		testCheckHashValid(t, "sha512")
	})

	t.Run("sha512_invalid", func(t *testing.T) {
		testCheckHashInvalid(t, "sha512")
	})

	t.Run("sha512_empty", func(t *testing.T) {
		testCheckHashEmpty(t, "sha512")
	})

	t.Run("sha3", func(t *testing.T) {
		testCheckHashValid(t, "sha3")
	})

	t.Run("sha3_invalid", func(t *testing.T) {
		testCheckHashInvalid(t, "sha3")
	})

	t.Run("sha3_empty", func(t *testing.T) {
		testCheckHashEmpty(t, "sha3")
	})

	t.Run("sha256", func(t *testing.T) {
		testCheckHashValid(t, "sha256")
	})

	t.Run("sha256_invalid", func(t *testing.T) {
		testCheckHashInvalid(t, "sha256")
	})

	t.Run("sha256_empty", func(t *testing.T) {
		testCheckHashEmpty(t, "sha256")
	})

	t.Run("sha1", func(t *testing.T) {
		testCheckHashValid(t, "sha1")
	})

	t.Run("sha1_invalid", func(t *testing.T) {
		testCheckHashInvalid(t, "sha1")
	})

	t.Run("sha1_empty", func(t *testing.T) {
		testCheckHashEmpty(t, "sha1")
	})

	t.Run("md5", func(t *testing.T) {
		testCheckHashValid(t, "md5")
	})

	t.Run("md5_invalid", func(t *testing.T) {
		testCheckHashInvalid(t, "md5")
	})

	t.Run("md5_empty", func(t *testing.T) {
		testCheckHashEmpty(t, "md5")
	})

	t.Run("no_media", func(t *testing.T) {
		t.Parallel()

		require := require.New(t)
		dir := t.TempDir()

		conf := medhash.Config{
			Dir: dir,
		}

		man := &medhash.Manifest{
			Version: medhash.ManifestFormatVer,
			Media:   []medhash.Media{},
			Config:  conf,
		}
		require.Error(man.Check("payload"))
	})
}

func testCheckHashInvalid(t *testing.T, alg string) {
	t.Parallel()

	require := require.New(t)
	conf, payload := testCheckHashCommon(t, alg, true, "__INVALID__")

	man := &medhash.Manifest{
		Version: medhash.ManifestFormatVer,
		Media:   []medhash.Media{payload},
		Config:  conf,
	}
	require.Error(man.Check(payload.Path))
}

func testCheckHashEmpty(t *testing.T, alg string) {
	t.Parallel()

	require := require.New(t)
	conf, payload := testCheckHashCommon(t, alg, true, "")

	man := &medhash.Manifest{
		Version: medhash.ManifestFormatVer,
		Media:   []medhash.Media{payload},
		Config:  conf,
	}
	require.NoError(man.Check(payload.Path))
}

func testCheckHashValid(t *testing.T, alg string) {
	t.Parallel()

	require := require.New(t)
	conf, payload := testCheckHashCommon(t, alg, false, "")

	man := &medhash.Manifest{
		Version: medhash.ManifestFormatVer,
		Media:   []medhash.Media{payload},
		Config:  conf,
	}
	require.NoError(man.Check(payload.Path))
}

func testCheckHashCommon(t *testing.T, alg string, valSet bool, val string) (medhash.Config, medhash.Media) {
	require := require.New(t)
	dir := t.TempDir()
	payload := testcommon.GenPayload(t, dir, testcommon.PayloadSize())
	require.FileExists(filepath.Join(dir, payload.Path))

	conf := medhash.Config{
		Dir: dir,
	}
	switch alg {
	case "xxh3":
		conf.XXH3 = true
		if valSet {
			payload.Hash.XXH3 = val
		}
	case "sha512":
		conf.SHA512 = true
		if valSet {
			payload.Hash.SHA512 = val
		}
	case "sha3":
		conf.SHA3 = true
		if valSet {
			payload.Hash.SHA3 = val
		}
	case "sha256":
		conf.SHA256 = true
		if valSet {
			payload.Hash.SHA256 = val
		}
	case "sha1":
		conf.SHA1 = true
		if valSet {
			payload.Hash.SHA1 = val
		}
	case "md5":
		conf.MD5 = true
		if valSet {
			payload.Hash.MD5 = val
		}
	}

	return conf, payload
}

func testGenHash(t *testing.T, alg string, assertFn func(t testing.TB, a *assert.Assertions,
	man *medhash.Manifest, pld medhash.Media)) {
	t.Parallel()

	require := require.New(t)
	dir := t.TempDir()
	payload := testcommon.GenPayload(t, dir, testcommon.PayloadSize())
	require.FileExists(filepath.Join(dir, payload.Path))

	conf := medhash.Config{
		Dir: dir,
	}
	switch alg {
	case "xxh3":
		conf.XXH3 = true
	case "sha512":
		conf.SHA512 = true
	case "sha3":
		conf.SHA3 = true
	case "sha256":
		conf.SHA256 = true
	case "sha1":
		conf.SHA1 = true
	case "md5":
		conf.MD5 = true
	}

	man, err := medhash.NewWithConfig(conf)
	require.NoError(err)
	require.NoError(man.Add(payload.Path))
	assertFn(t, assert.New(t), man, payload)
}
