package medhash_test

import (
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/stretchr/testify/suite"
)

const maxBuffer = 1 * 1024 * 1024 * 1024

type MedHashTestSuite struct {
	suite.Suite
}

func (m *MedHashTestSuite) GenPayload(t testing.TB, dir string, size int64) (payload string) {
	var buf []byte
	var counter int64

	t.Log("Generating payload")

	if size > maxBuffer {
		buf = make([]byte, maxBuffer)
	} else {
		buf = make([]byte, size)
	}

	payload = "payload"

	f, err := os.Create(filepath.Join(dir, payload))
	if err != nil {
		t.Fatalf("Cannot generate payload: %v", err)
		return
	}

	for counter < size {
		n, err := rand.Read(buf)
		if err != nil {
			t.Fatalf("Cannot generate payload: %v", err)
			return
		}

		if counter+int64(n) > int64(size) {
			n, err = f.Write(buf[:size-counter])
		} else {
			n, err = f.Write(buf)
		}
		if err != nil {
			t.Fatalf("Cannot generate payload: %v", err)
		}

		counter += int64(n)
	}

	err = f.Close()
	if err != nil {
		t.Fatalf("Cannot generate payload: %v", err)
	}

	t.Log("Done generating payload")

	return
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
