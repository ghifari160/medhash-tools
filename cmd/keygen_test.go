package cmd_test

import (
	"os"
	"path/filepath"

	"github.com/ghifari160/medhash-tools/cmd"
	"github.com/ghifari160/medhash-tools/testcommon"
)

func (s *CmdSuite) TestKeygen() {
	testKeyExists := func(s *CmdSuite, pregenPrivKey, pregenPubKey, force, success bool) {
		testcommon.SpoofDirectory(s.T())

		dir, err := cmd.ConfigDir()
		s.Require().NoError(err)

		if pregenPrivKey {
			s.Require().NoFileExists(filepath.Join(dir, "ed25519.key"))
			f, err := os.Create(filepath.Join(dir, "ed25519.key"))
			s.Require().NoError(err)
			s.Require().NoError(f.Close())
		}

		if pregenPubKey {
			s.Require().NoFileExists(filepath.Join(dir, "ed25519.pub"))
			f, err := os.Create(filepath.Join(dir, "ed25519.pub"))
			s.Require().NoError(err)
			s.Require().NoError(f.Close())
		}

		s.T().Cleanup(func() {
			err := os.RemoveAll(dir)
			s.Require().NoError(err)
		})

		c := new(cmd.Keygen)
		c.Algorithms = "ed25519"
		c.Force = force

		status := c.Execute()
		if success {
			s.Zero(status)
		} else {
			s.NotZero(status)
		}
	}

	s.Run("privKey_exists", func() {
		s.Run("forced", func() {
			testKeyExists(s, true, false, true, true)
		})

		s.Run("not_forced", func() {
			testKeyExists(s, true, false, false, false)
		})
	})

	s.Run("pubKey_exists", func() {
		s.Run("forced", func() {
			testKeyExists(s, false, true, true, true)
		})

		s.Run("not_forced", func() {
			testKeyExists(s, false, true, false, false)
		})
	})

	s.Run("unknown", func() {
		testcommon.SpoofDirectory(s.T())

		c := new(cmd.Keygen)
		c.Algorithms = "unknown"

		s.NotZero(c.Execute())
	})

	s.Run("ed25519", func() {
		testcommon.SpoofDirectory(s.T())

		dir, err := cmd.ConfigDir()
		s.Require().NoError(err)

		s.T().Cleanup(func() {
			err := os.RemoveAll(dir)
			s.Require().NoError(err)
		})

		c := new(cmd.Keygen)
		c.Algorithms = "ed25519"

		status := c.Execute()
		s.Require().Zero(status)

		s.Require().FileExists(filepath.Join(dir, "ed25519.key"))
		s.Require().FileExists(filepath.Join(dir, "ed25519.pub"))
	})
}
