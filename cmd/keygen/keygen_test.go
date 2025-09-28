package keygen_test

import (
	"context"
	"crypto/ed25519"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/ghifari160/medhash-tools/cmd/keygen"
	"github.com/ghifari160/medhash-tools/testcommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestKeygen(t *testing.T) {
	cases := []testcommon.TestCase{
		testcommon.Case("ed25519", "ed25519"),
		testcommon.Case("minisign/unencrypted", "minisign"),
		testcommon.Case("minisign/encrypted", "minisign",
			withPassword("correct-battery-horse-staple")),
	}

	testcommon.RunCases(t, testKeygen, cases)
}

func testKeygen(t *testing.T, alg string, opts ...testcommon.Options) {
	switch alg {
	case "ed25519":
		testEd25519(t)

	case "minisign":
		testMinisign(t, opts...)

	default:
		t.Fatal(errors.ErrUnsupported)
	}
}

func testEd25519(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	dir := t.TempDir()
	privPath := filepath.Join(dir, "private.key")
	pubPath := filepath.Join(dir, "public.key")

	var shouldError bool

	command := keygen.Command()
	command.ExitErrHandler = func(ctx context.Context, c *cli.Command, err error) {
		if shouldError {
			require.Error(err)
		} else {
			require.NoError(err)
		}
	}

	arguments := make([]string, 7)
	arguments[0] = "keygen"
	arguments[1] = "ed25519"
	arguments[2] = "--private"
	arguments[3] = privPath
	arguments[4] = "--public"
	arguments[5] = pubPath
	arguments[6] = "--force"

	r, w, err := os.Pipe()
	require.NoError(err)
	oldStdIn := os.Stdin
	t.Cleanup(func() {
		os.Stdin = oldStdIn
	})
	os.Stdin = r
	w.WriteString("\n")

	err = command.Run(t.Context(), arguments)
	require.NoError(err)

	require.FileExists(privPath)
	privKey := testcommon.LoadPEMKey(t, privPath)
	assert.Len(privKey, ed25519.PrivateKeySize)

	require.FileExists(pubPath)
	pubKey := testcommon.LoadPEMKey(t, pubPath)
	assert.Len(pubKey, ed25519.PublicKeySize)
}

func testMinisign(t *testing.T, opts ...testcommon.Options) {
	require := require.New(t)
	assert := assert.New(t)
	dir := t.TempDir()
	privPath := filepath.Join(dir, "private.key")
	pubPath := filepath.Join(dir, "public.key")

	options := testcommon.MergeOptions(opts...)
	password := options.Str("password")

	var shouldError bool

	command := keygen.Command()
	command.ExitErrHandler = exitHandler(t, &shouldError)

	arguments := make([]string, 7)
	arguments[0] = "keygen"
	arguments[1] = "minisign"
	arguments[2] = "--private"
	arguments[3] = privPath
	arguments[4] = "--public"
	arguments[5] = pubPath
	arguments[6] = "--force"

	r, w, err := os.Pipe()
	require.NoError(err)
	oldStdIn := os.Stdin
	t.Cleanup(func() {
		os.Stdin = oldStdIn
	})
	os.Stdin = r
	if password != "" {
		w.WriteString(password + "\n")
	} else {
		w.WriteString("\n")
	}

	err = command.Run(t.Context(), arguments)
	require.NoError(err)

	require.FileExists(privPath)
	privKey := testcommon.LoadMinisignPrivKey(t, privPath, password)

	require.FileExists(pubPath)
	pubKey := testcommon.LoadMinisignPubKey(t, pubPath)
	assert.True(pubKey.Equal(privKey.Public()))
}

func exitHandler(t testing.TB, shouldError *bool) cli.ExitErrHandlerFunc {
	t.Helper()
	require := require.New(t)
	return func(ctx context.Context, c *cli.Command, err error) {
		if *shouldError {
			require.Error(err)
		} else {
			require.NoError(err)
		}
	}
}

func withPassword(password string) testcommon.Options {
	return testcommon.NewOptions("password", password)
}
