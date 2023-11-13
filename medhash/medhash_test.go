package medhash_test

import (
	"encoding/json"
	"testing"

	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/ghifari160/medhash-tools/testcommon"
	"github.com/stretchr/testify/suite"
)

type MedHashTestSuite struct {
	suite.Suite

	PayloadSize int64
}

func (m *MedHashTestSuite) GenPayload(t testing.TB, dir string, size int64) (payload string) {
	return testcommon.GenPayload(t, dir, size).Path
}

func (s *MedHashTestSuite) TestNew() {
	s.Run("new", func() {
		var manifest *medhash.Manifest

		s.Require().NotPanics(func() {
			manifest = medhash.New()
		})

		s.NotNil(manifest)
		s.True(manifest.Config.XXH3)
	})

	s.Run("newWithConfig", func() {
		var manifest *medhash.Manifest

		config := medhash.Config{
			SHA3: true,
		}

		s.Require().NotPanics(func() {
			manifest = medhash.NewWithConfig(config)
		})

		s.NotNil(manifest)
		s.Equal(config, manifest.Config)
	})
}

func (s *MedHashTestSuite) TestManifest() {
	dir := s.T().TempDir()
	payload := testcommon.GenPayload(s.T(), dir, s.PayloadSize)

	manifest := medhash.New()
	manifest.Media = append(manifest.Media, payload)

	s.Run("TestJSON", func() {
		expect, err := json.MarshalIndent(manifest, "", "  ")
		s.Require().NoError(err)

		j, err := manifest.JSON()
		s.Require().NoError(err)

		s.Equal(expect, j)
	})

	s.Run("TestString", func() {
		expect, err := json.MarshalIndent(manifest, "", "  ")
		s.Require().NoError(err)

		var str string

		s.Require().NotPanics(func() {
			str = manifest.String()
		})

		s.Equal(string(expect), str)
	})

	s.Run("TestCopy", func() {
		signature := &medhash.Signature{
			Ed25519: "test_signature_hash",
		}
		manifest.Signature = signature

		s.Run("NoMedia", func() {
			target := medhash.New()
			m := manifest
			m.Media = make([]medhash.Media, 0)

			s.Require().NotPanics(func() {
				m.Copy(target)
			})

			s.Equal(m, target)
		})

		s.Run("NoSignature", func() {
			target := medhash.New()
			m := manifest
			m.Signature = nil

			s.Require().NotPanics(func() {
				m.Copy(target)
			})

			s.Equal(m, target)
		})

		s.Run("Full", func() {
			target := medhash.New()
			m := manifest

			s.Require().NotPanics(func() {
				m.Copy(target)
			})

			s.Equal(m, target)
		})
	})

	s.Run("TestStripSignature", func() {
		signature := &medhash.Signature{
			Ed25519: "test_signature_hash",
		}
		manifest.Signature = signature

		stripped, sig := manifest.StripSignature()

		s.Equal(manifest.Version, stripped.Version)
		s.Equal(manifest.Generator, stripped.Generator)
		s.Equal(manifest.Media, stripped.Media)
		s.Equal(manifest.Config, stripped.Config)
		s.Equal(signature, sig)
	})
}

func TestMedHash(t *testing.T) {
	s := new(MedHashTestSuite)

	if testing.Short() {
		s.PayloadSize = 1024
	} else {
		s.PayloadSize = 1 * 1024 * 1024 * 1024
	}

	suite.Run(t, s)
}
