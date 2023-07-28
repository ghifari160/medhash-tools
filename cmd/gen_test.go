package cmd_test

import (
	"os"
	"path/filepath"

	"github.com/ghifari160/medhash-tools/cmd"
	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/ghifari160/medhash-tools/testcommon"
)

func (s *CmdSuite) TestGen() {
	dir := s.T().TempDir()
	payload := testcommon.GenPayload(s.T(), dir, s.PayloadSize)

	s.Run("xxh3", func() {
		s.T().Cleanup(func() {
			err := os.Remove(filepath.Join(dir, medhash.DefaultManifestName))
			s.Require().NoError(err)
		})

		c := new(cmd.Gen)
		c.Dirs = []string{dir}
		c.Default = true
		c.XXH3 = true

		status := c.Execute()
		s.Require().Zero(status)

		s.Require().FileExists(filepath.Join(dir, medhash.DefaultManifestName))

		config := medhash.Config{
			XXH3: true,
		}

		testcommon.VerifyManifest(s.T(), dir, config, payload.Hash)
	})

	s.Run("sha3", func() {
		s.T().Cleanup(func() {
			err := os.Remove(filepath.Join(dir, medhash.DefaultManifestName))
			s.Require().NoError(err)
		})

		c := new(cmd.Gen)
		c.Dirs = []string{dir}
		c.Default = true
		c.SHA3 = true

		status := c.Execute()
		s.Require().Zero(status)

		s.Require().FileExists(filepath.Join(dir, medhash.DefaultManifestName))

		config := medhash.Config{
			SHA3: true,
		}

		testcommon.VerifyManifest(s.T(), dir, config, payload.Hash)
	})

	s.Run("sha256", func() {
		s.T().Cleanup(func() {
			err := os.Remove(filepath.Join(dir, medhash.DefaultManifestName))
			s.Require().NoError(err)
		})

		c := new(cmd.Gen)
		c.Dirs = []string{dir}
		c.Default = true
		c.SHA256 = true

		status := c.Execute()
		s.Require().Zero(status)

		s.Require().FileExists(filepath.Join(dir, medhash.DefaultManifestName))

		config := medhash.Config{
			SHA256: true,
		}

		testcommon.VerifyManifest(s.T(), dir, config, payload.Hash)
	})

	s.Run("sha1", func() {
		s.T().Cleanup(func() {
			err := os.Remove(filepath.Join(dir, medhash.DefaultManifestName))
			s.Require().NoError(err)
		})

		c := new(cmd.Gen)
		c.Dirs = []string{dir}
		c.Default = true
		c.SHA1 = true

		status := c.Execute()
		s.Require().Zero(status)

		s.Require().FileExists(filepath.Join(dir, medhash.DefaultManifestName))

		config := medhash.Config{
			SHA1: true,
		}

		testcommon.VerifyManifest(s.T(), dir, config, payload.Hash)
	})

	s.Run("md5", func() {
		s.T().Cleanup(func() {
			err := os.Remove(filepath.Join(dir, medhash.DefaultManifestName))
			s.Require().NoError(err)
		})

		c := new(cmd.Gen)
		c.Dirs = []string{dir}
		c.Default = true
		c.MD5 = true

		status := c.Execute()
		s.Require().Zero(status)

		s.Require().FileExists(filepath.Join(dir, medhash.DefaultManifestName))

		config := medhash.Config{
			MD5: true,
		}

		testcommon.VerifyManifest(s.T(), dir, config, payload.Hash)
	})

	s.Run("all", func() {
		s.T().Cleanup(func() {
			err := os.Remove(filepath.Join(dir, medhash.DefaultManifestName))
			s.Require().NoError(err)
		})

		c := new(cmd.Gen)
		c.Dirs = []string{dir}
		c.Default = true
		c.All = true

		status := c.Execute()
		s.Require().Zero(status)

		s.Require().FileExists(filepath.Join(dir, medhash.DefaultManifestName))

		testcommon.VerifyManifest(s.T(), dir, medhash.DefaultConfig, payload.Hash)
	})

	s.Run("default", func() {
		s.T().Cleanup(func() {
			err := os.Remove(filepath.Join(dir, medhash.DefaultManifestName))
			s.Require().NoError(err)
		})

		c := new(cmd.Gen)
		c.Dirs = []string{dir}
		c.Default = true

		status := c.Execute()
		s.Require().Zero(status)

		s.Require().FileExists(filepath.Join(dir, medhash.DefaultManifestName))

		testcommon.VerifyManifest(s.T(), dir, medhash.DefaultConfig, payload.Hash)
	})
}
