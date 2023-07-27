package cmd_test

import (
	"os"
	"path/filepath"

	"github.com/ghifari160/medhash-tools/cmd"
	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/ghifari160/medhash-tools/testcommon"
)

func (s *CmdSuite) TestUpgrade() {
	dir := s.T().TempDir()
	payload := testcommon.GenPayload(s.T(), dir, s.PayloadSize)

	s.Run("v0.1.0", func() {
		testcommon.CreateLegacyManifest(s.T(), dir, payload)
		s.T().Cleanup(func() {
			err := os.Remove(filepath.Join(dir, "sums.txt"))
			s.Require().NoError(err)
		})

		c := new(cmd.Upgrade)
		c.Dirs = []string{dir}
		c.Ignores = []string{"sums.txt"}
		c.Default = true

		status := c.Execute()
		s.Require().Zero(status)

		s.Require().FileExists(filepath.Join(dir, medhash.DefaultManifestName))

		testcommon.VerifyManifest(s.T(), dir, medhash.DefaultConfig, payload.Hash)
	})

	s.Run("v0.2.0", func() {
		testcommon.CreateManifest(s.T(), dir, payload, "0.2.0", medhash.Config{
			SHA256: true,
		})
		s.T().Cleanup(func() {
			err := os.Remove(filepath.Join(dir, medhash.DefaultManifestName))
			s.Require().NoError(err)
		})

		c := new(cmd.Upgrade)
		c.Default = true
		c.Dirs = []string{dir}

		status := c.Execute()
		s.Require().Zero(status)

		s.Require().FileExists(filepath.Join(dir, medhash.DefaultManifestName))

		testcommon.VerifyManifest(s.T(), dir, medhash.DefaultConfig, payload.Hash)
	})

	s.Run("v0.3.0", func() {
		testcommon.CreateManifest(s.T(), dir, payload, "0.3.0", medhash.Config{
			SHA256: true,
			SHA1:   true,
			MD5:    true,
		})
		s.T().Cleanup(func() {
			err := os.Remove(filepath.Join(dir, medhash.DefaultManifestName))
			s.Require().NoError(err)
		})

		c := new(cmd.Upgrade)
		c.Default = true
		c.Dirs = []string{dir}

		status := c.Execute()
		s.Require().Zero(status)

		s.Require().FileExists(filepath.Join(dir, medhash.DefaultManifestName))

		testcommon.VerifyManifest(s.T(), dir, medhash.DefaultConfig, payload.Hash)
	})

	s.Run("v0.4.0", func() {
		s.Run("not_forced", func() {
			testcommon.CreateManifest(s.T(), dir, payload, "0.4.0", medhash.DefaultConfig)
			s.T().Cleanup(func() {
				err := os.Remove(filepath.Join(dir, medhash.DefaultManifestName))
				s.Require().NoError(err)
			})

			c := new(cmd.Upgrade)
			c.Default = true
			c.Dirs = []string{dir}

			status := c.Execute()
			s.Require().NotZero(status)
		})

		s.Run("forced", func() {
			testcommon.CreateManifest(s.T(), dir, payload, "0.4.0", medhash.DefaultConfig)
			s.T().Cleanup(func() {
				err := os.Remove(filepath.Join(dir, medhash.DefaultManifestName))
				s.Require().NoError(err)
			})

			c := new(cmd.Upgrade)
			c.Force = true
			c.Default = true
			c.Dirs = []string{dir}

			status := c.Execute()
			s.Require().Zero(status)

			testcommon.VerifyManifest(s.T(), dir, medhash.DefaultConfig, payload.Hash)
		})
	})
}
