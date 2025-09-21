package cmd

import (
	"encoding/pem"
	"fmt"
	"os"
)

// Generator generates a new keypair for a given algorithm.
type Generator func() (pubKey, privKey []byte, err error)

// Storer stores key for a given algorithm.
type Storer func(private bool, path string, data []byte) error

// Loader loads PEM encoded key in path, asserting the type matches expectedType.
func Loader(path, expectedType string) (key []byte, err error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return
	}

	var block *pem.Block
	data := f
	for block == nil && data != nil && len(data) > 0 {
		block, data = pem.Decode(data)
		if block != nil && block.Type == expectedType {
			key = block.Bytes
			return
		}
	}
	if block == nil || block.Type != expectedType {
		err = fmt.Errorf("no valid key in %s", path)
	}
	return
}
