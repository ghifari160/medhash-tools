package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ghifari160/medhash-tools/src/data"
	"github.com/ghifari160/medhash-tools/src/packageinfo"
)

// Deprecated: legacy code.
func handleUpgradeV010(root string) error {
	legacyMedhashFile, err := ioutil.ReadFile(root + "/sums.txt")

	if err != nil {
		return err
	}

	legacyMedhash := string(legacyMedhashFile)
	legacyMedhashMedia := strings.Split(legacyMedhash, "\n")

	var media []data.Media

	fmt.Println("Upgrading legacy medhash file")

	errCount := 0
	for i := 0; i < len(legacyMedhashMedia); i++ {
		legacyMedia := strings.Split(legacyMedhashMedia[i], "  ")

		if len(legacyMedia) > 1 {
			legacyPath := legacyMedia[1]
			legacyHash := legacyMedia[0]

			if len(legacyPath) > 2 && legacyPath[0:2] == "./" {
				legacyPath = legacyPath[2:]
			}

			compPath := legacyPath

			if root != "." {
				compPath = root + "/" + legacyPath
			}

			fmt.Print("  ")
			fmt.Print(compPath)
			fmt.Print(": ")

			var hash data.Hash
			hash, err = data.GenHash(compPath)

			if err != nil {
				errCount++
				fmt.Println("Error")
			} else {
				if hash.SHA256 != legacyHash {
					errCount++
					fmt.Println("Error")
				} else {
					fmt.Println("OK")

					media = append(media,
						data.Media{
							Path: legacyPath,
							Hash: hash,
						})
				}
			}
		}
	}

	if errCount > 0 {
		fmt.Fprintln(os.Stderr, "Error: Media integrity error detected!")
		return errors.New("media integrity")
	}

	fmt.Println("Sanity checking upgraded medhash")

	errCount = 0
	for i := 0; i < len(media); i++ {
		berr := false

		compPath := media[i].Path

		if root != "." {
			compPath = root + "/" + compPath
		}

		fmt.Print("  ")
		fmt.Print(compPath)
		fmt.Print(": ")

		berr = !data.ChkHash(compPath, media[i].Hash)

		if berr {
			errCount++
			fmt.Println("Error")
		} else {
			fmt.Println("OK")
		}
	}

	if errCount > 0 {
		fmt.Fprintln(os.Stderr, "Error: Media integrity error detected!")
		return errors.New("media integrity")
	}

	medhash := data.Medhash{
		Version: packageinfo.Version,
		Media:   media,
	}

	medhashJSON, err := json.MarshalIndent(medhash, "", "    ")
	if err != nil {
		return err
	}

	medhashFile, err := os.Create(root + "/medhash.json")
	if err != nil {
		return err
	}

	_, err = medhashFile.Write(medhashJSON)
	if err != nil {
		return err
	}

	return nil
}
