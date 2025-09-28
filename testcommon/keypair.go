package testcommon

import (
	"crypto/ed25519"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"aead.dev/minisign"
	"github.com/stretchr/testify/require"
)

// LoadPEMKey loads a PEM encoded key from path.
func LoadPEMKey(t testing.TB, path string) (key []byte) {
	t.Helper()

	require := require.New(t)

	f, err := os.ReadFile(path)
	require.NoError(err)

	block, _ := pem.Decode(f)
	require.NotNil(block)
	return block.Bytes
}

// GemEd25519Keypair generates a new Ed25519 keypair and returns the paths.
// The generated keypair are stored in a temporary directory associated with t.
func GenEd25519Keypair(t testing.TB) (pubPath, privPath string) {
	t.Helper()

	require := require.New(t)

	t.Log("Generating Ed25519 keypair")
	dir := t.TempDir()

	pubKey, privKey, err := ed25519.GenerateKey(nil)
	require.NoError(err)

	privPath = filepath.Join(dir, "ed25519.key")
	pubPath = filepath.Join(dir, "ed25519.pub")

	require.NoError(storePEMKey(privPath, "ED25519 PRIVATE KEY", privKey))
	require.NoError(storePEMKey(pubPath, "ED25519 PUBLIC KEY", pubKey))

	return
}

func LoadMinisignPrivKey(t testing.TB, path, password string) (key minisign.PrivateKey) {
	t.Helper()

	require := require.New(t)

	f, err := os.ReadFile(path)
	require.NoError(err)

	err = key.UnmarshalText(f)
	if err != nil {
		key, err = minisign.DecryptKey(password, f)
		require.NoError(err)
	}
	return
}

func LoadMinisignPubKey(t testing.TB, path string) (key minisign.PublicKey) {
	t.Helper()

	require := require.New(t)

	f, err := os.ReadFile(path)
	require.NoError(err)

	require.NoError(key.UnmarshalText(f))

	return
}

// GenMinisignKeypair generates a new Minisign keypair and returns the paths.
// The generated keypair are stored in a temporary directory associated with t.
// If password is an empty string, the key will be stored unencrypted.
// Otherwise, GenMinisignKeypair encrypts the private key with a key derived from password.
func GenMinisignKeypair(t testing.TB, password string) (pubPath, privPath string) {
	t.Helper()

	require := require.New(t)

	t.Log("Generating Minisign keypair")
	dir := t.TempDir()

	pubKey, privKey, err := minisign.GenerateKey(nil)
	require.NoError(err)

	privPath = filepath.Join(dir, "minisign.key")
	var privKeyText []byte
	if password != "" {
		privKeyText, err = minisign.EncryptKey(password, privKey)
	} else {
		privKeyText, err = privKey.MarshalText()
	}
	require.NoError(err)
	require.NoError(os.WriteFile(privPath, privKeyText, 0600))

	pubPath = filepath.Join(dir, "minisign.pub")
	pubKeyText, err := pubKey.MarshalText()
	require.NoError(err)
	require.NoError(os.WriteFile(pubPath, pubKeyText, 0600))

	return
}

// storePEMKey PEM encodes key in path.
func storePEMKey(path, keyType string, key []byte) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	block := pem.Block{
		Type:  keyType,
		Bytes: key,
	}

	return pem.Encode(f, &block)
}
