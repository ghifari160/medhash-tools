package medhash

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"

	"aead.dev/minisign"
)

// StripSignature returns a copy of man with the Signature field set to the zero value.
func (man *Manifest) StripSignature() *Manifest {
	stripped := new(Manifest)
	*stripped = *man
	stripped.Signature = Signature{}
	return stripped
}

// Sign signs man.
func (man *Manifest) Sign() error {
	errs := make([]error, 0)
	if man.Config.Ed25519.Enabled {
		errs = append(errs, ed25519_sign(man))
	}
	if man.Config.Minisign.Enabled {
		errs = append(errs, minisign_sign(man))
	}
	return errors.Join(errs...)
}

// Verify verifies all signatures in man.
func (man *Manifest) Verify() error {
	errs := make([]error, 0)
	if man.Config.Ed25519.Enabled {
		errs = append(errs, ed25519_verify(man))
	}
	if man.Config.Minisign.Enabled {
		errs = append(errs, minisign_verify(man))
	}
	return errors.Join(errs...)
}

// ed25519_sign signs man with Ed25519.
func ed25519_sign(man *Manifest) (err error) {
	payload, err := man.StripSignature().JSON()
	if err != nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			err = recoverErr(r)
		}
	}()

	signature := ed25519.Sign(man.Config.Ed25519.PrivKey, payload)
	man.Signature.Ed25519 = hex.EncodeToString(signature)

	return
}

// ed25519_verify verifies man with Ed25519.
func ed25519_verify(man *Manifest) (err error) {
	rawSig, err := hex.DecodeString(man.Signature.Ed25519)
	if err != nil {
		return fmt.Errorf("ed25519: malformed signature: %w", err)
	}

	payload, err := man.StripSignature().JSON()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			err = recoverErr(r)
		}
	}()

	if !ed25519.Verify(man.Config.Ed25519.PubKey, payload, rawSig) {
		err = fmt.Errorf("ed25519: bad signature")
	}

	return
}

// minisign_sign signs manifest with Minisign.
func minisign_sign(manifest *Manifest) (err error) {
	var buf bytes.Buffer

	err = manifest.StripSignature().JSONStream(&buf)
	if err != nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			err = recoverErr(r)
		}
	}()

	reader := minisign.NewReader(&buf)
	_, err = io.Copy(io.Discard, reader)

	manifestName := manifest.Config.Manifest
	if manifestName == "" {
		manifestName = DefaultManifestName
	}
	trusted := fmt.Sprintf("timestamp:%d\tfile:%s", time.Now().Unix(), manifestName)
	untrusted := fmt.Sprintf("signature of stripped %s (key id: %d)",
		manifestName, manifest.Config.Minisign.PrivKey.ID())

	sig := reader.SignWithComments(manifest.Config.Minisign.PrivKey, trusted, untrusted)
	manifest.Signature.Minisign = string(sig)

	return
}

// minisign_verify verifies manifest with Minisign.
func minisign_verify(manifest *Manifest) (err error) {
	var buf bytes.Buffer

	err = manifest.StripSignature().JSONStream(&buf)
	if err != nil {
		return
	}

	rawSig := []byte(manifest.Signature.Minisign)

	reader := minisign.NewReader(&buf)
	_, err = io.Copy(io.Discard, reader)
	if err != nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			err = recoverErr(r)
		}
	}()

	if !reader.Verify(manifest.Config.Minisign.PubKey, rawSig) {
		err = fmt.Errorf("minisign: bad signature")
	}

	return
}

// recoverErr recovers error from r.
// If r is an error or can be casted into the error interface, it is casted and returned.
// Otherwise, recoverErr returns fmt.Errorf("%v", r).
//
// Example:
//
//	func willPanic() error {
//		var err error
//		defer func() {
//			if r := recover(); r != nil {
//				err = recoverErr(r)
//			}
//		}()
//		somePanicCode()
//		return err
//	}
func recoverErr(r any) error {
	if r == nil {
		return nil
	}

	if e, ok := r.(error); ok {
		return e
	} else {
		return fmt.Errorf("%v", r)
	}
}

// Signature stores signatures for the Manifest.
type Signature struct {
	Ed25519  string `json:"ed25519,omitempty"`
	Minisign string `json:"minisign,omitempty"`
}

// SigConf configures a given signature algorithm.
type SigConf struct {
	// Enable signing or verification.
	Enabled bool
	// PubKey stores the unencoded public key.
	// PubKey affects only verification.
	PubKey []byte
	// PrivKey stores the unencoded private key.
	// PrivKey affects only signing.
	PrivKey []byte
}

type GenericSigConf[S any, P any] struct {
	// Enable signing or verification.
	Enabled bool
	// PubKey stores the unencoded public key.
	// PubKey affects only verification.
	PubKey P
	// PrivKey stores the unencoded private key.
	// PrivKey affects only signing.
	PrivKey S
}
