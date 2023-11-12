package medhash

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
)

const (
	PRIVTYPE = "PRIVATE KEY"
	PUBTYPE  = "PUBLIC KEY"
)

var (
	ErrNoSig        = errors.New("no signature")
	ErrMalformedSig = errors.New("malformed signature")
	ErrBadPrivKey   = errors.New("bad private key")
	ErrBadPublicKey = errors.New("bad public key")
)

// GenKey generates an Ed25519 key-pair for signing purposes.
func GenKey() (pubKey, privKey []byte, err error) {
	return ed25519.GenerateKey(nil)
}

// EncodeKey encodes a key into PEM-encoded format.
func EncodeKey(key []byte, private bool) []byte {
	b := new(pem.Block)
	b.Bytes = key

	if private {
		b.Type = PRIVTYPE
	} else {
		b.Type = PUBTYPE
	}

	return pem.EncodeToMemory(b)
}

// DecodeKey decodes a key from PEM-encoded format.
func DecodeKey(encoded []byte) (key []byte, private bool) {
	b, _ := pem.Decode(encoded)
	if b == nil {
		return
	}

	if b == nil {
		return
	}

	key = b.Bytes
	private = b.Type == PRIVTYPE

	return
}

// Sign signs a Manifest.
// Only enabled algorithms are generated.
func Sign(config Config, manifest *Manifest) (signed *Manifest, err error) {
	stripped, _ := manifest.StripSignature()
	sig := new(Signature)

	j, err := stripped.JSON()
	if err != nil {
		return
	}

	if config.Ed25519.Enabled {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("ed25519: %w: %v", ErrBadPrivKey, r)
			}
		}()

		s := ed25519.Sign(config.Ed25519.PrivateKey, j)
		sig.Ed25519 = hex.EncodeToString(s)
	}

	stripped.Signature = sig
	signed = stripped

	return
}

// Verify verifies the Signature of manifest.
// Only enabled algorithms are verified.
// If an algorithm is enabled and its corresponding Signature field is empty, Verify returns
// ErrNoSig.
func Verify(config Config, manifest *Manifest) (valid bool, err error) {
	stripped, sig := manifest.StripSignature()

	j, err := stripped.JSON()
	if err != nil {
		return
	}

	var rawSig []byte

	if config.Ed25519.Enabled {
		if len(sig.Ed25519) < 1 {
			err = fmt.Errorf("ed25519: %w", ErrNoSig)
			return
		}

		rawSig, err = hex.DecodeString(sig.Ed25519)
		if err != nil {
			err = fmt.Errorf("ed25519: %w: %w", ErrMalformedSig, err)
			return
		}

		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("ed25519: %w: %v", ErrBadPublicKey, r)
			}
		}()

		if !ed25519.Verify(config.Ed25519.PublicKey, j, rawSig) {
			return
		}
	}

	valid = true
	return
}
