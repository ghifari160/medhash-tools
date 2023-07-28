package medhash_test

import "github.com/ghifari160/medhash-tools/medhash"

func (s *MedHashTestSuite) TestGenHash() {
	dir := s.T().TempDir()

	payload := s.GenPayload(s.T(), dir, s.PayloadSize)

	s.Run("xxh3", func() {
		conf := medhash.Config{
			Dir:  dir,
			Path: payload,
			XXH3: true,
		}

		m, err := medhash.GenHash(conf)
		s.Require().NoError(err)

		s.NotEmpty(m.Hash.XXH3)
	})

	s.Run("sha3", func() {
		conf := medhash.Config{
			Dir:  dir,
			Path: payload,
			SHA3: true,
		}

		m, err := medhash.GenHash(conf)
		s.Require().NoError(err)

		s.NotEmpty(m.Hash.SHA3_256)
	})

	s.Run("sha256", func() {
		conf := medhash.Config{
			Dir:    dir,
			Path:   payload,
			SHA256: true,
		}

		m, err := medhash.GenHash(conf)
		s.Require().NoError(err)

		s.NotEmpty(m.Hash.SHA256)
	})

	s.Run("sha1", func() {
		conf := medhash.Config{
			Dir:  dir,
			Path: payload,
			SHA1: true,
		}

		m, err := medhash.GenHash(conf)
		s.Require().NoError(err)

		s.NotEmpty(m.Hash.SHA1)
	})

	s.Run("md5", func() {
		conf := medhash.Config{
			Dir:  dir,
			Path: payload,
			MD5:  true,
		}

		m, err := medhash.GenHash(conf)
		s.Require().NoError(err)

		s.NotEmpty(m.Hash.MD5)
	})

	s.Run("all", func() {
		conf := medhash.Config{
			Dir:  dir,
			Path: payload,

			SHA3:   true,
			SHA256: true,
			SHA1:   true,
			MD5:    true,
		}

		m, err := medhash.GenHash(conf)
		s.Require().NoError(err)

		s.NotEmpty(m.Hash.SHA3_256)
		s.NotEmpty(m.Hash.SHA256)
		s.NotEmpty(m.Hash.SHA1)
		s.NotEmpty(m.Hash.MD5)
	})
}
