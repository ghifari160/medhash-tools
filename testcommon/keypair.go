package testcommon

import (
	"crypto/ed25519"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// LoadKey loads a PEM encoded key from path.
func LoadKey(t testing.TB, path string) (key []byte) {
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

	require.NoError(storeKey(privPath, "ED25519 PRIVATE KEY", privKey))
	require.NoError(storeKey(pubPath, "ED25519 PUBLIC KEY", pubKey))

	return
}

// storeKey PEM encodes key in path.
func storeKey(path, keyType string, key []byte) error {
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
