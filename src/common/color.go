// MedHash Tools
// Copyright (c) 2021 GHIFARI160
// MIT License

// Deprecated: legacy code.
package common

import (
	"fmt"
	"os"
)

// Deprecated: legacy code.
func IsColorTerm() bool {
	return os.Getenv("TERM") != ""
}

// Deprecated: legacy code.
func ColorPrintln(asciiCode string, a ...interface{}) {
	if IsColorTerm() {
		fmt.Print(asciiCode)
		fmt.Print(a...)
		fmt.Println("\x1B[0m")
	} else {
		fmt.Println(a...)
	}
}
