package testcommon

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// SpoofDirectory spoofs user directories (home, cache, config, etc.) to a new temporary directory
// for testing purposes.
// The directory is automatically removed by t.Cleanup.
// SpoofDirectory affects the whole process, so it cannot be used in parallel tests.
func SpoofDirectory(t testing.TB) {
	t.Helper()

	require := require.New(t)

	dir := t.TempDir()

	switch runtime.GOOS {
	case "windows":
		t.Setenv("LocalAppData", filepath.Join(dir, "LocalAppData"))
		t.Setenv("AppData", filepath.Join(dir, "AppData"))
		t.Setenv("USERPROFILE", dir)

	case "darwin", "ios":
		t.Setenv("HOME", dir)

	case "plan9":
		t.Setenv("home", dir)

	default:
		t.Setenv("XDG_CACHE_HOME", filepath.Join(dir, ".cache"))
		t.Setenv("XDG_CONFIG_HOME", filepath.Join(dir, ".config"))
		t.Setenv("HOME", dir)
	}

	d, err := os.UserCacheDir()
	require.NoError(err)
	err = os.MkdirAll(d, 0755)
	require.NoError(err)

	d, err = os.UserConfigDir()
	require.NoError(err)
	err = os.MkdirAll(d, 0755)
	require.NoError(err)

	d, err = os.UserHomeDir()
	require.NoError(err)
	err = os.MkdirAll(d, 0755)
	require.NoError(err)
}
