// MedHash Tools
// Copyright (c) 2021 GHIFARI160
// MIT License

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ghifari160/medhash-tools/src/common"
	"github.com/ghifari160/medhash-tools/src/medhash"
)

const NAME string = "medhash-gen"

const DEFAULT_FLAGMANIFEST string = "__TARGET__"

func main() {
	targetDir := "."

	var flagVersion bool
	flag.BoolVar(&flagVersion, "version", false, "Print version")

	var flagVerbose bool
	flag.BoolVar(&flagVerbose, "v", false, "Verbose mode")

	var flagManifest string
	flag.StringVar(&flagManifest, "manifest", DEFAULT_FLAGMANIFEST, "Manifest output path")

	flag.Parse()

	if len(flag.Args()) > 0 {
		targetDir = path.Clean(flag.Args()[0])
	}

	common.PrintHeader(NAME)

	if flagVersion {
		os.Exit(0)
	}

	if flagManifest == DEFAULT_FLAGMANIFEST {
		flagManifest = path.Join(targetDir, medhash.MEDHASH_MANIFEST_NAME)
	}

	homeDir, _ := os.UserHomeDir()

	if strings.HasPrefix(targetDir, "~") {
		targetDir = path.Join(homeDir, targetDir[1:])
	}

	if strings.HasPrefix(flagManifest, "~") {
		flagManifest = path.Join(homeDir, flagManifest[1:])
	}

	cwd, _ := os.Getwd()

	if flagVerbose {
		fmt.Printf("Working Dir: %s\n", cwd)
		fmt.Printf("Target Dir: %s\n", targetDir)
		fmt.Printf("Manifest: %s\n", flagManifest)
	}

	fmt.Println("Generating hash files")

	var files []string
	var infos []os.FileInfo

	err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		infos = append(infos, info)
		return nil
	})
	common.HandleError(err, 1)

	var media []medhash.Media
	mediaMap := make(map[string]*medhash.Hash)

	var manifestIgnorePath string
	if targetDir != "." {
		manifestIgnorePath = path.Join(targetDir, medhash.MEDHASH_MANIFEST_NAME)
	} else {
		manifestIgnorePath = medhash.MEDHASH_MANIFEST_NAME
	}

	errMap := make(map[string]error)

	for i := 0; i < len(files); i++ {
		if infos[i].Mode().IsRegular() && files[i] != manifestIgnorePath {
			fmt.Printf("  %s\n", files[i])

			hash, err := medhash.GenHash(files[i])
			if err == nil {
				mediaMap[files[i]] = hash
			} else {
				errMap[files[i]] = err
			}
		}
	}

	fmt.Println("Sanity checking files")

	sanitizedMediaMap := make(map[string]*medhash.Hash)
	invalidCount := 0

	for k, v := range mediaMap {
		fmt.Printf("  %s: ", k)

		valid, err := medhash.ChkHash(k, v)
		if err != nil {
			errMap[k] = err
		}

		if !valid {
			invalidCount++
			common.ColorPrintln("\x1B[31m", "ERROR")
		} else {
			common.ColorPrintln("\x1B[32m", "OK")

			relPath, err := filepath.Rel(targetDir, k)
			if err != nil {
				errMap[k] = err
			} else {
				sanitizedMediaMap[relPath] = v
			}
		}
	}

	media = make([]medhash.Media, len(sanitizedMediaMap))
	mI := 0
	for k, v := range sanitizedMediaMap {
		media[mI] = medhash.Media{
			Path: k,
			Hash: v,
		}

		mI++
	}

	if len(errMap) > 0 {
		for p, e := range errMap {
			fmt.Fprintf(os.Stderr, "  %s: %v\n", p, e)
		}
	}

	if invalidCount > 0 {
		fmt.Fprintln(os.Stderr, "Media integrity error detected!")
	}

	medHash := medhash.New()
	medHash.Media = media

	medhashJSON, err := json.MarshalIndent(medHash, "", "    ")
	common.HandleError(err, 1)

	medhashFile, err := os.Create(flagManifest)
	common.HandleError(err, 1)

	_, err = medhashFile.Write(medhashJSON)
	common.HandleError(err, 1)

	err = medhashFile.Sync()
	common.HandleError(err, 1)

	fmt.Println("Done!")
}
