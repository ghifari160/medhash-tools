// MedHash Tools
// Copyright (c) 2021 GHIFARI160
// MIT License

package common

import "fmt"

func PrintHeader(toolName string) {
	fmt.Printf("%s v%s\n", NAME, VERSION)
	fmt.Printf("%s\n", toolName)
}
