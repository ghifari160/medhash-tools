package main

import "errors"

// Deprecated: legacy code.
func handleUpgradeV030(root string) error {
	return errors.New("medhash file version is current")
}
