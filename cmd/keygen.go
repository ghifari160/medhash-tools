package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghifari160/medhash-tools/color"
	"github.com/ghifari160/medhash-tools/medhash"
)

// Keygen command generate keypairs for Manifest signing.
type Keygen struct {
	Algorithms string `arg:"positional" help:"key algorithm"`
	Force      bool   `arg:"--force" help:"force regeneration of keys"`

	Ed25519 bool `arg:"-"`
}

func (k *Keygen) parse() error {
	for _, alg := range strings.Split(k.Algorithms, ",") {
		switch strings.ToLower(alg) {
		case "ed25519":
			k.Ed25519 = true

		default:
			return fmt.Errorf("unknown algorithm: %s", alg)
		}
	}

	return nil
}

func (k *Keygen) Execute() (status int) {
	err := k.parse()
	if err != nil {
		color.Printf("%s %v\n", PrefixError, err)
		status = 1

		return
	}

	dir, err := ConfigDir()
	if err != nil {
		color.Printf("%s %v\n", PrefixError, err)
		status = 1

		return
	}

	pubKeys := make(map[string][]byte)

	if k.Ed25519 {
		pubEncoded, err := k.ed25519(dir)
		if err != nil {
			color.Printf("%s %v\n", PrefixError, err)
			status = 1

			return
		}

		pubKeys["Ed25519"] = pubEncoded
	}

	color.Println()

	for alg, key := range pubKeys {
		color.Printf("%s Public Key:\n%s%s%s\n", alg, color.Green, key, color.Reset)
	}

	return
}

// ed25519 generates ed25519 keypair for Manifest signing.
// It stores the private key to `ConfigDir()/ed25519.key` in PEM-encoded format, and the public key
// to `./ed25519.pub` in PEM-encoded format.
func (k *Keygen) ed25519(dir string) (pubEncoded []byte, err error) {
	color.Println("Generating Ed25519 keypair")
	pub, priv, err := medhash.GenKey()
	if err != nil {
		return
	}

	pubEncoded = medhash.EncodeKey(pub, false)
	privEncoded := medhash.EncodeKey(priv, true)

	err = k.writeKey(privEncoded, filepath.Join(dir, "ed25519.key"), k.Force, true)
	if err != nil {
		return
	}

	err = k.writeKey(pubEncoded, filepath.Join(dir, "ed25519.pub"), k.Force, false)
	if err != nil {
		return
	}

	return
}

// writeKey writes key to path.
// If a file exists at path and overwrite is true, the file is overwritten.
// Otherwise, an error is returned.
func (k *Keygen) writeKey(key []byte, path string, overwrite bool, private bool) (err error) {
	flag := os.O_CREATE | os.O_WRONLY
	if !overwrite {
		flag |= os.O_EXCL
	}

	var perm fs.FileMode
	if private {
		perm = 0600
	} else {
		perm = 0644
	}

	var keyType string
	if private {
		keyType = "Private"
	} else {
		keyType = "Public"
	}

	color.Printf("  Writing %s key to %s\n", keyType, path)
	f, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.Write(key)
	if err != nil {
		return
	}

	return
}
