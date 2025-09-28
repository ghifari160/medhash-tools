package medhash_test

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"io"
	"path/filepath"
	"testing"

	"aead.dev/minisign"
	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/ghifari160/medhash-tools/testcommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSign(t *testing.T) {
	t.Parallel()

	cases := []testcommon.TestCase{
		testcommon.Case("ed25519/valid", "ed25519"),
		testcommon.Case("ed25519/invalid_key_length", "ed25519", withInvalidKey(true)),

		testcommon.Case("minisign/valid", "minisign"),
	}

	testcommon.RunCases(t, testSign, cases)
}

func TestVerify(t *testing.T) {
	t.Parallel()

	cases := []testcommon.TestCase{
		testcommon.Case("ed25519/valid", "ed25519"),
		testcommon.Case("ed25519/invalid_key_length", "ed25519", withInvalidKey(false)),
		testcommon.Case("ed25519/malformed_signature", "ed25519", withMalformedSignature()),
		testcommon.Case("ed25519/bad_signature", "ed25519", withBadSignature()),

		testcommon.Case("minisign/valid", "minisign"),
		testcommon.Case("minisign/malformed_signature", "minisign", withMalformedSignature()),
		testcommon.Case("minisign/bad_signature", "minisign", withBadSignature()),
	}

	testcommon.RunCases(t, testVerify, cases)
}

func testSign(t *testing.T, alg string, opts ...testcommon.Options) {
	switch alg {
	case "ed25519":
		ed25519_testSign(t, opts...)

	case "minisign":
		minisign_testSign(t, opts...)
	}
}

func testVerify(t *testing.T, alg string, opts ...testcommon.Options) {
	switch alg {
	case "ed25519":
		ed25519_testVerify(t, opts...)

	case "minisign":
		minisign_testVerify(t, opts...)
	}
}

func ed25519_testSign(t *testing.T, opts ...testcommon.Options) {
	t.Parallel()

	require := require.New(t)
	assert := assert.New(t)
	dir := t.TempDir()
	payload := testcommon.GenPayload(t, dir, testcommon.PayloadSize())
	require.FileExists(filepath.Join(dir, payload.Path))
	pubKey, privKey := ed25519_genKey(t)

	options := testcommon.MergeOptions(opts...)
	invalidPrivateKey := options.Bool("invalid_private_key")

	if invalidPrivateKey {
		privKey = privKey[:len(privKey)/2]
	}

	conf := medhash.Config{
		Dir: dir,
		Ed25519: medhash.SigConf{
			Enabled: true,
			PubKey:  pubKey,
			PrivKey: privKey,
		},
	}

	man := medhash.Manifest{
		Version: medhash.ManifestFormatVer,
		Media:   []medhash.Media{payload},
		Config:  conf,
	}

	err := man.Sign()
	if invalidPrivateKey {
		require.Error(err)
	} else {
		require.NoError(err)
	}

	rawSig, err := hex.DecodeString(man.Signature.Ed25519)
	require.NoError(err)

	if !invalidPrivateKey {
		j, err := man.StripSignature().JSON()
		require.NoError(err)
		var valid bool
		require.NotPanics(func() {
			valid = ed25519.Verify(conf.Ed25519.PubKey, j, rawSig)
		})
		assert.True(valid)
	}
}

func ed25519_testVerify(t *testing.T, opts ...testcommon.Options) {
	t.Parallel()

	require := require.New(t)
	dir := t.TempDir()
	payload := testcommon.GenPayload(t, dir, testcommon.PayloadSize())
	require.FileExists(filepath.Join(dir, payload.Path))
	pubKey, privKey := ed25519_genKey(t)

	options := testcommon.MergeOptions(opts...)
	invalidPublicKey := options.Bool("invalid_public_key")
	malformedSignature := options.Bool("malformed_signature")
	badSignature := options.Bool("bad_signature")

	if invalidPublicKey {
		pubKey = pubKey[:len(pubKey)/2]
	}

	conf := medhash.Config{
		Dir: dir,
		Ed25519: medhash.SigConf{
			Enabled: true,
			PubKey:  pubKey,
			PrivKey: privKey,
		},
	}

	man := medhash.Manifest{
		Version: medhash.ManifestFormatVer,
		Media:   []medhash.Media{payload},
		Config:  conf,
	}

	j, err := man.StripSignature().JSON()
	require.NoError(err)
	var sig []byte
	require.NotPanics(func() {
		sig = ed25519.Sign(conf.Ed25519.PrivKey, j)
	})

	if badSignature {
		man.Signature.Ed25519 = hex.EncodeToString([]byte("bad_signature"))
	} else if malformedSignature {
		man.Signature.Ed25519 = "malformed_signature"
	} else {
		man.Signature.Ed25519 = hex.EncodeToString(sig)
	}

	err = man.Verify()
	if invalidPublicKey || malformedSignature || badSignature {
		require.Error(err)
	} else {
		require.NoError(err)
	}
}

func ed25519_genKey(t testing.TB) (pubKey, privKey []byte) {
	t.Helper()
	require := require.New(t)
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	require.NoError(err)
	return
}

func minisign_testSign(t *testing.T, _ ...testcommon.Options) {
	t.Parallel()

	require := require.New(t)
	assert := assert.New(t)
	dir := t.TempDir()
	payload := testcommon.GenPayload(t, dir, testcommon.PayloadSize())
	require.FileExists(filepath.Join(dir, payload.Path))
	pubKey, privKey := minisign_genKey(t)

	conf := medhash.Config{
		Dir:      dir,
		Manifest: medhash.DefaultManifestName,
		Minisign: medhash.GenericSigConf[minisign.PrivateKey, minisign.PublicKey]{
			Enabled: true,
			PubKey:  pubKey,
			PrivKey: privKey,
		},
	}

	man := medhash.Manifest{
		Version: medhash.ManifestFormatVer,
		Media:   []medhash.Media{payload},
		Config:  conf,
	}

	err := man.Sign()
	require.NoError(err)

	var buf bytes.Buffer

	rawSig := []byte(man.Signature.Minisign)
	require.NoError(man.StripSignature().JSONStream(&buf))

	reader := minisign.NewReader(&buf)
	_, err = io.Copy(io.Discard, reader)
	require.NoError(err)

	var valid bool
	require.NotPanics(func() {
		valid = reader.Verify(pubKey, rawSig)
	})
	assert.True(valid)
}

func minisign_testVerify(t *testing.T, opts ...testcommon.Options) {
	t.Parallel()

	require := require.New(t)
	dir := t.TempDir()
	payload := testcommon.GenPayload(t, dir, testcommon.PayloadSize())
	require.FileExists(filepath.Join(dir, payload.Path))
	pubKey, privKey := minisign_genKey(t)

	options := testcommon.MergeOptions(opts...)
	malformedSignature := options.Bool("malformed_signature")
	badSignature := options.Bool("bad_signature")

	conf := medhash.Config{
		Dir:      dir,
		Manifest: medhash.DefaultManifestName,
		Minisign: medhash.GenericSigConf[minisign.PrivateKey, minisign.PublicKey]{
			Enabled: true,
			PubKey:  pubKey,
			PrivKey: privKey,
		},
	}

	man := medhash.Manifest{
		Version: medhash.DefaultManifestName,
		Media:   []medhash.Media{payload},
		Config:  conf,
	}

	var buf bytes.Buffer

	err := man.StripSignature().JSONStream(&buf)
	require.NoError(err)
	reader := minisign.NewReader(&buf)
	_, err = io.Copy(io.Discard, reader)
	require.NoError(err)

	var shouldError bool
	if badSignature {
		shouldError = true
		man.Signature.Minisign = base64.StdEncoding.EncodeToString([]byte("bad_signature"))
	} else if malformedSignature {
		shouldError = true
		man.Signature.Minisign = "malformed_signature"
	} else {
		require.NotPanics(func() {
			sig := reader.Sign(privKey)
			man.Signature.Minisign = string(sig)
		})
	}

	err = man.Verify()
	if shouldError {
		require.Error(err)
	} else {
		require.NoError(err)
	}
}

func minisign_genKey(t testing.TB) (minisign.PublicKey, minisign.PrivateKey) {
	t.Helper()
	require := require.New(t)
	pubKey, privKey, err := minisign.GenerateKey(nil)
	require.NoError(err)
	return pubKey, privKey
}
