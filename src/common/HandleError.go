// MedHash Tools
// Copyright (c) 2021 GHIFARI160
// MIT License

package common

import (
	"fmt"
	"os"
)

// Deprecated: legacy code.
func HandleError(err error, exitCode int) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(exitCode)
	}
}
