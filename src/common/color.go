// MedHash Tools
// Copyright (c) 2021 GHIFARI160
// MIT License

package common

import (
	"fmt"
	"os"
)

func IsColorTerm() bool {
	return os.Getenv("TERM") != ""
}

func ColorPrintln(asciiCode string, a ...interface{}) {
	if IsColorTerm() {
		fmt.Print(asciiCode)
		fmt.Print(a...)
		fmt.Println("\x1B[0m")
	} else {
		fmt.Println(a...)
	}
}
