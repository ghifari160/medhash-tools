package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghifari160/medhash-tools/src/data"
)

// Deprecated: legacy code.
func handleUpgradeV020(root string) error {
	medhashLegacyFile, err := ioutil.ReadFile(root + "/medhash.json")

	if err != nil {
		return err
	}

	fmt.Println("Upgrading medhash file")

	var medhash data.Medhash
	err = json.Unmarshal(medhashLegacyFile, &medhash)

	medhash.Version = "0.3.0"

	fmt.Println("Sanity checking media")

	errCount := 0
	for i := 0; i < len(medhash.Media); i++ {
		compPath := medhash.Media[i].Path

		if root != "." {
			compPath = root + "/" + compPath
		}

		fmt.Print("  ")
		fmt.Print(compPath)
		fmt.Print(": ")

		berr := !data.ChkHash(compPath, medhash.Media[i].Hash)

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
