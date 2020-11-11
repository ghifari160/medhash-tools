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

	medhashFile, err := ioutil.ReadFile(root + "/medhash.json")

	if err != nil {
		panic(err)
	}

	var medhash data.Medhash
	err = json.Unmarshal(medhashFile, &medhash)

	fmt.Println("Checking hash files")

	errCount := 0
	for i := 0; i < len(medhash.Media); i++ {
		berr := false

		path := ""
		if root != "." {
			path = root + "/" + medhash.Media[i].Path
		} else {
			path = medhash.Media[i].Path
		}

		fmt.Print("  ")
		fmt.Print(path)
		fmt.Print(": ")

		berr = !data.ChkHash(path, medhash.Media[i].Hash)

		if berr {
			errCount++
			fmt.Println("Error")
		} else {
			fmt.Println("OK")
		}
	}

	if errCount > 0 {
		println("Media integrity error detected!")
	}

	fmt.Println("Done!")
}
