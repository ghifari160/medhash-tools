package keygen_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/ghifari160/medhash-tools/cmd/keygen"
	"github.com/ghifari160/medhash-tools/testcommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestKeygen(t *testing.T) {
	t.Parallel()

	cases := []testcommon.TestCase{
		testcommon.Case("ed25519", "ed25519"),
	}

	testcommon.RunCases(t, testKeygen, cases)
}

func testKeygen(t *testing.T, alg string, opts ...testcommon.Options) {
	t.Parallel()

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
	arguments[2] = "--private"
	arguments[3] = privPath
	arguments[4] = "--public"
	arguments[5] = pubPath
	arguments[6] = "--force"

	var privLen, pubLen int

	switch alg {
	case "ed25519":
		arguments[1] = "ed25519"
		privLen = 64
		pubLen = 32
	default:
		t.Fatal(errors.ErrUnsupported)
	}

	err := command.Run(t.Context(), arguments)
	require.NoError(err)

	require.FileExists(privPath)
	privKey := testcommon.LoadKey(t, privPath)
	assert.Len(privKey, privLen)

	require.FileExists(pubPath)
	pubKey := testcommon.LoadKey(t, pubPath)
	assert.Len(pubKey, pubLen)
}
