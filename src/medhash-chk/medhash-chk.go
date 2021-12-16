// MedHash Tools
// Copyright (c) 2021 GHIFARI160
// MIT License

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/ghifari160/medhash-tools/src/common"
	"github.com/ghifari160/medhash-tools/src/medhash"
)

const NAME string = "medhash-chk"

func main() {
	targetDir := "."

	var flagVersion bool
	flag.BoolVar(&flagVersion, "version", false, "Print version")

	var flagVerbose bool
	flag.BoolVar(&flagVerbose, "v", false, "Verbose mode")

	flag.Parse()

	if len(flag.Args()) > 0 {
		targetDir = flag.Args()[0]
	}

	common.PrintHeader(NAME)

	if flagVersion {
		os.Exit(0)
	}

	cwd, _ := os.Getwd()

	homeDir, _ := os.UserHomeDir()

	if strings.HasPrefix(targetDir, "~") {
		targetDir = path.Join(homeDir, targetDir[1:])
	}

	if flagVerbose {
		fmt.Printf("Working Dir: %s\n", cwd)
		fmt.Printf("Target Dir: %s\n", targetDir)
	}

	medhashFile, err := ioutil.ReadFile(path.Join(targetDir, medhash.MEDHASH_MANIFEST_NAME))
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s does not exists in the target directory\n", medhash.MEDHASH_MANIFEST_NAME)
			os.Exit(1)
		} else {
			common.HandleError(err, 1)
		}
	}

	var medHash medhash.MedHash
	err = json.Unmarshal(medhashFile, &medHash)

	fmt.Println("Checking hash files")

	invalidCount := 0
	errMap := make(map[string]error)

	for i := 0; i < len(medHash.Media); i++ {
		mediaPath := ""
		if targetDir != "." {
			mediaPath = path.Join(targetDir, medHash.Media[i].Path)
		} else {
			mediaPath = medHash.Media[i].Path
		}

		fmt.Printf("  %s: ", mediaPath)

		valid, err := medhash.ChkHash(mediaPath, medHash.Media[i].Hash)
		if err != nil {
			errMap[medHash.Media[i].Path] = err
		}

		if !valid {
			invalidCount++
			fmt.Println("Error")
		} else {
			fmt.Println("OK")
		}
	}

	if len(errMap) > 0 {
		for p, e := range errMap {
			fmt.Fprintf(os.Stderr, "  %s: %v\n", p, e)
		}
	}

	if invalidCount > 0 {
		fmt.Fprintln(os.Stderr, "Media integrity error detected!")
	}

	fmt.Println("Done!")
}
