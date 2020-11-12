package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghifari160/medhash-tools/src/data"
	"github.com/ghifari160/medhash-tools/src/packageinfo"
)

func main() {
	root := "."

	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	fmt.Print(packageinfo.Name)
	fmt.Print(" v")
	fmt.Println(packageinfo.Version)

	cwd, _ := os.Getwd()

	fmt.Print("Working Dir: ")
	fmt.Println(cwd)

	fmt.Print("Target Dir: ")
	fmt.Println(root)

	_, err := os.Stat(root + "/medhash.json")

	if err != nil && os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "medhash.json not found")

		_, err = os.Stat(root + "/sums.txt")

		if err != nil && os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "Error: medhash file does not exist in this directory!")
			os.Exit(-1)
		} else {
			fmt.Println("Legacy medhash file (v0.1.0) found")

			err := handleUpgradeV010(root)

			if err != nil {
				fmt.Fprintln(os.Stderr, "Error: Unable to upgrade legacy medhash file")
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(-1)
			} else {
				fmt.Println("Done!")
				os.Exit(0)
			}
		}
	}

	medhashFile, err := ioutil.ReadFile(root + "/medhash.json")

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: Unable to read medhash file")
		os.Exit(-1)
	}

	var medhash data.Medhash
	err = json.Unmarshal(medhashFile, &medhash)

	fmt.Print("medhash.json found: medhash v")
	fmt.Println(medhash.Version)

	if medhash.Version == "0.2.0" {
		err = handleUpgradeV020(root)
	} else if medhash.Version == "0.3.0" {
		err = handleUpgradeV030(root)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: Unable to upgrade medhash file")
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}

	fmt.Println("Done!")
}
