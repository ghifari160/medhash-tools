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
	"path/filepath"
	"strings"

	"github.com/ghifari160/medhash-tools/src/common"
	"github.com/ghifari160/medhash-tools/src/medhash"
)

const NAME string = "medhash-chk"

const DEFAULT_FLAGFILE string = "__FILE__"
const DEFAULT_FLAGMANIFEST string = "__TARGET__"

func main() {
	targetDir := "."

	var flagVersion bool
	flag.BoolVar(&flagVersion, "version", false, "Print version")

	var flagVerbose bool
	flag.BoolVar(&flagVerbose, "v", false, "Verbose mode")

	var flagFile string
	flag.StringVar(&flagFile, "file", DEFAULT_FLAGFILE, "Verify a specific file in the manifest")

	var flagManifest string
	flag.StringVar(&flagManifest, "manifest", DEFAULT_FLAGMANIFEST, "Path to Manifest")

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

	if flagManifest == DEFAULT_FLAGMANIFEST {
		flagManifest = filepath.Join(targetDir, medhash.MEDHASH_MANIFEST_NAME)
	}

	if strings.HasPrefix(flagManifest, "~") {
		flagManifest = filepath.Join(homeDir, flagManifest[1:])
	}

	if strings.HasPrefix(targetDir, "~") {
		targetDir = filepath.Join(homeDir, targetDir[1:])
	}

	if flagVerbose {
		fmt.Printf("Working Dir: %s\n", cwd)
		fmt.Printf("Target Dir: %s\n", targetDir)
		fmt.Printf("Manifest: %s\n", flagManifest)
	}

	medhashFile, err := ioutil.ReadFile(flagManifest)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Manifest does not exists in the target directory")
			os.Exit(1)
		} else {
			common.HandleError(err, 1)
		}
	}

	var medHash medhash.MedHash
	err = json.Unmarshal(medhashFile, &medHash)
	common.HandleError(err, 1)

	if flagVerbose && medHash.Generator != "" {
		fmt.Printf("Generator: %s\n", medHash.Generator)
	}

	fmt.Println("Checking hash files")

	invalidCount := 0
	errMap := make(map[string]error)

	for i := 0; i < len(medHash.Media); i++ {
		skip := false

		if flagFile != DEFAULT_FLAGFILE && medHash.Media[i].Path != flagFile {
			skip = true
		}

		if !skip {
			mediaPath := ""
			if targetDir != "." {
				mediaPath = filepath.Join(targetDir, medHash.Media[i].Path)
			} else {
				mediaPath = medHash.Media[i].Path
			}

			fmt.Printf("  %s: ", mediaPath)

			valid, err := medhash.ChkHash(&medHash.Media[i])
			if err != nil {
				errMap[medHash.Media[i].Path] = err
			}

			if !valid {
				invalidCount++
				common.ColorPrintln("\x1B[31m", "ERROR")
			} else {
				common.ColorPrintln("\x1B[32m", "OK")
			}
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
