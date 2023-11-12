package medhash_test

import (
	"crypto/ed25519"
	"encoding/pem"

	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/ghifari160/medhash-tools/testcommon"
)

func (s *MedHashTestSuite) TestGenKey() {
	s.T().Parallel()

	pubKey, privKey, err := medhash.GenKey()
	s.Require().NoError(err)

	s.Require().NotNil(pubKey)
	s.Require().NotNil(privKey)

	s.Len(privKey, ed25519.PrivateKeySize)
	s.Len(pubKey, ed25519.PublicKeySize)
}

func (s *MedHashTestSuite) TestEncodeDecodeKey() {
	s.T().Parallel()

	pubKey, privKey, err := medhash.GenKey()
	s.Require().NoError(err)

	testEncodeKey := func(s *MedHashTestSuite, pemEncoded *[]byte, key []byte, private bool) {
		*pemEncoded = medhash.EncodeKey(key, private)

		s.Require().NotZero(len(*pemEncoded))

		b, rest := pem.Decode(*pemEncoded)
		s.Require().NotNil(b)
		s.Require().Empty(rest)
	}

	testDecodeKey := func(s *MedHashTestSuite, pemEncoded *[]byte, key []byte, private bool) {
		k, priv := medhash.DecodeKey(*pemEncoded)

		s.Equal(private, priv)
		s.Equal(key, k)
	}

	s.Run("pubKey", func() {
		var pemEncoded []byte

		s.Run("TestEncodeKey", func() {
			testEncodeKey(s, &pemEncoded, pubKey, false)
		})

		s.Run("TestDecodeKey", func() {
			testDecodeKey(s, &pemEncoded, pubKey, false)
		})
	})

	s.Run("privKey", func() {
		var pemEncoded []byte

		s.Run("TestEncodeKey", func() {
			testEncodeKey(s, &pemEncoded, privKey, true)
		})

		s.Run("TestDecodeKey", func() {
			testDecodeKey(s, &pemEncoded, privKey, true)
		})
	})
}

func (s *MedHashTestSuite) TestSignVerify() {
	dir := s.T().TempDir()
	payload := testcommon.GenPayload(s.T(), dir, s.PayloadSize)

	conf := medhash.DefaultConfig
	conf.Dir = dir

	manifest := medhash.New()
	manifest.Media = append(manifest.Media, payload)

	type testConf struct {
		signErr   bool
		signValid bool

		c medhash.Config
	}

	testSign := func(s *MedHashTestSuite, conf testConf,
		manifest *medhash.Manifest) (signed *medhash.Manifest) {
		signed, err := medhash.Sign(conf.c, manifest)
		if !conf.signErr {
			s.Require().NoError(err)
		} else {
			s.Require().Error(err)
		}

		if !conf.signErr {
			stripped, sig := signed.StripSignature()

			s.Equal(manifest, stripped)

			if conf.c.Ed25519.Enabled {
				s.Len(sig.Ed25519, ed25519.SignatureSize*2)
			}
		}

		return
	}

	testVerify := func(s *MedHashTestSuite, conf testConf, signed *medhash.Manifest) {
		valid, err := medhash.Verify(conf.c, signed)
		if !conf.signErr {
			s.Require().NoError(err)
		} else {
			s.Require().Error(err)
		}

		if !conf.signErr {
			s.Equal(conf.signValid, valid)
		}
	}

	s.Run("ed25519", func() {
		pub, priv, err := medhash.GenKey()
		s.Require().NoError(err)

		man := manifest
		conf := testConf{
			signErr:   false,
			signValid: true,

			c: conf,
		}

		conf.c.Ed25519.Enabled = true

		s.Run("good", func() {
			var signed *medhash.Manifest

			conf.signErr = false
			conf.signValid = true

			conf.c.Ed25519.PrivateKey = priv
			conf.c.Ed25519.PublicKey = pub

			s.Run("TestSign", func() {
				signed = testSign(s, conf, man)
			})

			s.Run("TestVerify", func() {
				testVerify(s, conf, signed)
			})
		})

		s.Run("failVerify", func() {
			conf.signErr = false
			conf.signValid = false

			conf.c.Ed25519.PublicKey = pub

			signed := medhash.New()
			man.Copy(signed)
			signed.Signature = &medhash.Signature{
				Ed25519: "aabbcc",
			}

			s.Run("TestVerify", func() {
				testVerify(s, conf, signed)
			})
		})

		s.Run("malformedSig", func() {
			conf.signErr = true

			conf.c.Ed25519.PublicKey = pub

			signed := medhash.New()
			man.Copy(signed)
			signed.Signature = &medhash.Signature{
				Ed25519: "invalid_signature",
			}

			s.Run("TestVerify", func() {
				testVerify(s, conf, signed)
			})
		})

		s.Run("badPrivKey", func() {
			conf.signErr = true

			conf.c.Ed25519.PrivateKey = []byte{0xAA, 0xBB, 0xCC}

			s.Run("TestSign", func() {
				testSign(s, conf, man)
			})
		})

		s.Run("badPubKey", func() {
			var signed *medhash.Manifest

			conf.signErr = false
			conf.signValid = true

			conf.c.Ed25519.PrivateKey = priv
			conf.c.Ed25519.PublicKey = []byte{0xAA, 0xBB, 0xCC}

			s.Run("TestSign", func() {
				signed = testSign(s, conf, man)
			})

			conf.signErr = true

			s.Run("TestVerify", func() {
				testVerify(s, conf, signed)
			})
		})
	})
}
