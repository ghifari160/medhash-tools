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

func TestMedHash(t *testing.T) {
	s := new(MedHashTestSuite)

	if testing.Short() {
		s.PayloadSize = 1024
	} else {
		s.PayloadSize = 1 * 1024 * 1024 * 1024
	}

	suite.Run(t, s)
}
