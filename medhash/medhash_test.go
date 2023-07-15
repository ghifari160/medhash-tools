// MedHash Tools
// Copyright (c) 2023 GHIFARI160
// MIT License

package medhash_test

import (
	"testing"

	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/ghifari160/medhash-tools/testcommon"
	"github.com/stretchr/testify/suite"
)

type MedHashTestSuite struct {
	suite.Suite
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
		s.True(manifest.Config.SHA3)
		s.True(manifest.Config.SHA256)
		s.True(manifest.Config.SHA1)
		s.True(manifest.Config.MD5)
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

func TestMedHash(t *testing.T) {
	suite.Run(t, new(MedHashTestSuite))
}
