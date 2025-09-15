package cmd_test

import (
	"os"
	"path/filepath"

	"github.com/ghifari160/medhash-tools/cmd"
	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/ghifari160/medhash-tools/testcommon"
	"github.com/stretchr/objx"
)

func (s *CmdSuite) TestChk() {
	dir := s.T().TempDir()
	payload := testcommon.GenPayload(s.T(), dir, s.PayloadSize)

	s.Run("xxh3", func() {
		testcommon.CreateManifest(s.T(), dir, payload, medhash.ManifestFormatVer, medhash.Config{
			XXH3: true,
		})
		s.T().Cleanup(func() {
			err := os.Remove(filepath.Join(dir, medhash.DefaultManifestName))
			s.Require().NoError(err)
		})

		c := new(cmd.Chk)
		c.Dirs = []string{dir}
		c.Default = true
		c.XXH3 = true

		status := c.Execute()
		s.Require().Zero(status)

		config := medhash.Config{
			XXH3: true,
		}

		testcommon.VerifyManifest(s.T(), dir, config, payload.Hash)
	})

	s.Run("sha512", func() {
		testcommon.CreateManifest(s.T(), dir, payload, medhash.ManifestFormatVer, medhash.Config{
			SHA512: true,
		})
		s.T().Cleanup(func() {
			err := os.Remove(filepath.Join(dir, medhash.DefaultManifestName))
			s.Require().NoError(err)
		})

		c := new(cmd.Chk)
		c.Dirs = []string{dir}
		c.Default = true
		c.SHA512 = true

		status := c.Execute()
		s.Require().Zero(status)

		config := medhash.Config{
			SHA512: true,
		}

		testcommon.VerifyManifest(s.T(), dir, config, payload.Hash)
	})

	s.Run("sha3", func() {
		testcommon.CreateManifest(s.T(), dir, payload, medhash.ManifestFormatVer, medhash.Config{
			SHA3: true,
		})
		s.T().Cleanup(func() {
			err := os.Remove(filepath.Join(dir, medhash.DefaultManifestName))
			s.Require().NoError(err)
		})

		c := new(cmd.Chk)
		c.Dirs = []string{dir}
		c.Default = true
		c.SHA3 = true

		status := c.Execute()
		s.Require().Zero(status)

		config := medhash.Config{
			SHA3: true,
		}

		testcommon.VerifyManifest(s.T(), dir, config, payload.Hash)
	})

	s.Run("sha256", func() {
		testcommon.CreateManifest(s.T(), dir, payload, medhash.ManifestFormatVer, medhash.Config{
			SHA256: true,
		})
		s.T().Cleanup(func() {
			err := os.Remove(filepath.Join(dir, medhash.DefaultManifestName))
			s.Require().NoError(err)
		})

		c := new(cmd.Chk)
		c.Dirs = []string{dir}
		c.Default = true
		c.SHA256 = true

		status := c.Execute()
		s.Require().Zero(status)

		config := medhash.Config{
			SHA256: true,
		}

		testcommon.VerifyManifest(s.T(), dir, config, payload.Hash)
	})

	s.Run("sha1", func() {
		testcommon.CreateManifest(s.T(), dir, payload, medhash.ManifestFormatVer, medhash.Config{
			SHA1: true,
		})
		s.T().Cleanup(func() {
			err := os.Remove(filepath.Join(dir, medhash.DefaultManifestName))
			s.Require().NoError(err)
		})

		c := new(cmd.Chk)
		c.Dirs = []string{dir}
		c.Default = true
		c.SHA1 = true

		status := c.Execute()
		s.Require().Zero(status)

		config := medhash.Config{
			SHA1: true,
		}

		testcommon.VerifyManifest(s.T(), dir, config, payload.Hash)
	})

	s.Run("md5", func() {
		testcommon.CreateManifest(s.T(), dir, payload, medhash.ManifestFormatVer, medhash.Config{
			MD5: true,
		})
		s.T().Cleanup(func() {
			err := os.Remove(filepath.Join(dir, medhash.DefaultManifestName))
			s.Require().NoError(err)
		})

		c := new(cmd.Chk)
		c.Dirs = []string{dir}
		c.Default = true
		c.MD5 = true

		status := c.Execute()
		s.Require().Zero(status)

		config := medhash.Config{
			MD5: true,
		}

		testcommon.VerifyManifest(s.T(), dir, config, payload.Hash)
	})

	s.Run("all", func() {
		testcommon.CreateManifest(s.T(), dir, payload, medhash.ManifestFormatVer,
			medhash.AllConfig)
		s.T().Cleanup(func() {
			err := os.Remove(filepath.Join(dir, medhash.DefaultManifestName))
			s.Require().NoError(err)
		})

		c := new(cmd.Chk)
		c.Dirs = []string{dir}
		c.Default = true
		c.All = true

		status := c.Execute()
		s.Require().Zero(status)

		testcommon.VerifyManifest(s.T(), dir, medhash.AllConfig, payload.Hash)
	})

	s.Run("default", func() {
		testcommon.CreateManifest(s.T(), dir, payload, medhash.ManifestFormatVer,
			medhash.DefaultConfig)
		s.T().Cleanup(func() {
			err := os.Remove(filepath.Join(dir, medhash.DefaultManifestName))
			s.Require().NoError(err)
		})

		c := new(cmd.Chk)
		c.Dirs = []string{dir}
		c.Default = true

		status := c.Execute()
		s.Require().Zero(status)

		testcommon.VerifyManifest(s.T(), dir, medhash.DefaultConfig, payload.Hash)
	})

	s.Run("manifest_path", func() {
		testcommon.CreateManifest(s.T(), dir, payload, medhash.ManifestFormatVer,
			medhash.DefaultConfig)
		err := os.Rename(filepath.Join(dir, medhash.DefaultManifestName),
			filepath.Join(dir, "manifest.json"))
		s.Require().NoError(err)
		s.T().Cleanup(func() {
			err := os.Remove(filepath.Join(dir, "manifest.json"))
			s.Require().NoError(err)
		})

		c := new(cmd.Chk)
		c.Dirs = []string{dir}
		c.Default = true
		c.Manifest = filepath.Join(dir, "manifest.json")

		status := c.Execute()
		s.Require().Zero(status)

		config := medhash.DefaultConfig

		manFile, err := os.ReadFile(filepath.Join(dir, "manifest.json"))
		s.Require().NoError(err)

		manifest, err := objx.FromJSON(string(manFile))
		s.Require().NoError(err)

		s.Require().Equal(medhash.ManifestFormatVer, manifest.Get("version").Str())

		if config.SHA512 {
			s.Equal(payload.Hash.SHA512, manifest.Get("media[0].hash.sha512").Str())
		}

		if config.SHA3 {
			s.Equal(payload.Hash.SHA3, manifest.Get("media[0].hash.sha3").Str())
			s.Equal(payload.Hash.SHA3_256, manifest.Get("media[0].hash.sha3-256").Str())
		}

		if config.SHA256 {
			s.Equal(payload.Hash.SHA256, manifest.Get("media[0].hash.sha256").Str())
		}

		if config.SHA1 {
			s.Equal(payload.Hash.SHA1, manifest.Get("media[0].hash.sha1").Str())
		}

		if config.MD5 {
			s.Equal(payload.Hash.MD5, manifest.Get("media[0].hash.md5").Str())
		}
	})
}
