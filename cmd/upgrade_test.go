package cmd_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ghifari160/medhash-tools/cmd"
	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/ghifari160/medhash-tools/testcommon"
	"github.com/stretchr/objx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *CmdSuite) TestUpgrade() {
	dir := s.T().TempDir()
	payload := testcommon.GenPayload(s.T(), dir, 1*1024*1024*1024)

	createLegacyManifest := func(t testing.TB, dir string, payload medhash.Media) {
		t.Helper()

		require := require.New(t)

		f, err := os.Create(filepath.Join(dir, "sums.txt"))
		require.NoError(err)

		_, err = f.WriteString(fmt.Sprintf("%s        %s\n",
			payload.Path, payload.Hash.SHA256))
		require.NoError(err)

		err = f.Close()
		require.NoError(err)
	}

	createManifest := func(t testing.TB, dir string, payload medhash.Media, ver string,
		config medhash.Config) {
		t.Helper()

		require := require.New(t)

		med := payload

		if !config.SHA3 {
			med.Hash.SHA3_256 = ""
		}

		if !config.SHA256 {
			med.Hash.SHA256 = ""
		}

		if !config.SHA1 {
			med.Hash.SHA1 = ""
		}

		if !config.MD5 {
			med.Hash.MD5 = ""
		}

		man := medhash.New()
		man.Version = ver
		man.Generator = "MedHash Tools Test"

		man.Media = []medhash.Media{med}

		manFile, err := json.Marshal(man)
		require.NoError(err)

		f, err := os.Create(filepath.Join(dir, medhash.DefaultManifestName))
		require.NoError(err)

		_, err = f.Write(manFile)
		require.NoError(err)

		err = f.Close()
		require.NoError(err)
	}

	verifyManifest := func(t testing.TB, dir string, config medhash.Config, hash medhash.Hash) {
		t.Helper()

		require := require.New(t)
		assert := assert.New(t)

		manFile, err := os.ReadFile(filepath.Join(dir, medhash.DefaultManifestName))
		require.NoError(err)

		manifest, err := objx.FromJSON(string(manFile))
		require.NoError(err)

		require.Equal(medhash.ManifestFormatVer, manifest.Get("version").Str())

		if config.SHA3 {
			assert.Equal(hash.SHA3_256, manifest.Get("media[0].hash.sha3-256").Str())
		}

		if config.SHA256 {
			assert.Equal(hash.SHA256, manifest.Get("media[0].hash.sha256").Str())
		}

		if config.SHA1 {
			assert.Equal(hash.SHA1, manifest.Get("media[0].hash.sha1").Str())
		}

		if config.MD5 {
			assert.Equal(hash.MD5, manifest.Get("media[0].hash.md5").Str())
		}
	}

	s.Run("v0.1.0", func() {
		createLegacyManifest(s.T(), dir, payload)
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

		config := medhash.Config{
			SHA3:   true,
			SHA256: true,
			SHA1:   true,
			MD5:    true,
		}

		verifyManifest(s.T(), dir, config, payload.Hash)
	})

	s.Run("v0.2.0", func() {
		createManifest(s.T(), dir, payload, "0.2.0", medhash.Config{
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

		config := medhash.Config{
			SHA3:   true,
			SHA256: true,
			SHA1:   true,
			MD5:    true,
		}

		verifyManifest(s.T(), dir, config, payload.Hash)
	})

	s.Run("v0.3.0", func() {
		createManifest(s.T(), dir, payload, "0.3.0", medhash.Config{
			SHA3:   true,
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

		config := medhash.Config{
			SHA3:   true,
			SHA256: true,
			SHA1:   true,
			MD5:    true,
		}

		verifyManifest(s.T(), dir, config, payload.Hash)
	})

	s.Run("v0.4.0", func() {
		s.Run("not_forced", func() {
			createManifest(s.T(), dir, payload, "0.4.0", medhash.Config{
				SHA3:   true,
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
			s.Require().NotZero(status)
		})

		s.Run("forced", func() {
			createManifest(s.T(), dir, payload, "0.4.0", medhash.Config{
				SHA3:   true,
				SHA256: true,
				SHA1:   true,
				MD5:    true,
			})
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

			config := medhash.Config{
				SHA3:   true,
				SHA256: true,
				SHA1:   true,
				MD5:    true,
			}

			verifyManifest(s.T(), dir, config, payload.Hash)
		})
	})
}
