package medhash_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/stretchr/testify/require"
)

func (s *MedHashTestSuite) TestChkHash() {
	dir := s.T().TempDir()

	payload := s.GenPayload(s.T(), dir, s.PayloadSize)

	config := medhash.Config{
		Dir:  dir,
		Path: payload,
	}

	s.Run("sha3", func() {
		config := config
		config.SHA3 = true

		med := genManifest(s.T(), config)

		config.Dir = dir

		valid, err := medhash.ChkHash(config, med)
		s.Require().NoError(err)

		s.True(valid)
	})

	s.Run("sha256", func() {
		config := config
		config.SHA256 = true

		med := genManifest(s.T(), config)

		config.Dir = dir

		valid, err := medhash.ChkHash(config, med)
		s.Require().NoError(err)

		s.True(valid)
	})

	s.Run("sha1", func() {
		config := config
		config.SHA1 = true

		med := genManifest(s.T(), config)

		config.Dir = dir

		valid, err := medhash.ChkHash(config, med)
		s.Require().NoError(err)

		s.True(valid)
	})

	s.Run("md5", func() {
		config := config
		config.MD5 = true

		med := genManifest(s.T(), config)

		config.Dir = dir

		valid, err := medhash.ChkHash(config, med)
		s.Require().NoError(err)

		s.True(valid)
	})

	s.Run("all", func() {
		config := config
		config.SHA3 = true
		config.SHA256 = true
		config.SHA1 = true
		config.MD5 = true

		med := genManifest(s.T(), config)

		config.Dir = dir

		valid, err := medhash.ChkHash(config, med)
		s.Require().NoError(err)

		s.True(valid)
	})

	s.Run("invalid", func() {
		config := config
		config.SHA3 = true
		config.SHA256 = true
		config.SHA1 = true
		config.MD5 = true

		med := genManifest(s.T(), config)
		med.Hash.SHA3_256 = "afafaf"

		config.Dir = dir

		valid, err := medhash.ChkHash(config, med)
		s.Require().NoError(err)

		s.False(valid)
	})
}

func genManifest(t testing.TB, config medhash.Config) medhash.Media {
	t.Helper()

	require := require.New(t)

	med, err := medhash.GenHash(config)
	require.NoError(err)

	manifest := medhash.New()
	manifest.Generator = "MedHash Tools Test"
	manifest.Media = []medhash.Media{med}

	manFile, err := json.Marshal(manifest)
	require.NoError(err)

	f, err := os.Create(filepath.Join(config.Dir, medhash.DefaultManifestName))
	require.NoError(err)
	defer f.Close()

	_, err = f.Write(manFile)
	require.NoError(err)

	t.Cleanup(func() {
		os.RemoveAll(filepath.Join(config.Dir, medhash.DefaultManifestName))
	})

	return med
}
