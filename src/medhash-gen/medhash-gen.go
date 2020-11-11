package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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

	fmt.Println("Generating hash files")

	var files []string
	var infos []os.FileInfo

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		infos = append(infos, info)
		return nil
	})

	if err != nil {
		panic(err)
	}

	var media []data.Media

	for i := 0; i < len(files); i++ {
		if root != "." && files[i] != cwd+"/medhash.json" && infos[i].Mode().IsRegular() {
			fmt.Print("  ")
			fmt.Println(files[i])

			hash, err := data.GenHash(files[i])
			if err == nil {
				media = append(media,
					data.Media{
						Path: files[i],
						Hash: hash,
					})
			}
		} else if root == "." && files[i] != "medhash.json" && files[i] != cwd+"/medhash.json" && infos[i].Mode().IsRegular() {
			fmt.Print("  ")
			fmt.Println(files[i])

			hash, err := data.GenHash(files[i])
			if err == nil {
				media = append(media,
					data.Media{
						Path: files[i],
						Hash: hash,
					})
			}
		}
	}

	fmt.Println("Sanity checking files")

	errCount := 0
	for i := 0; i < len(media); i++ {
		berr := false

		fmt.Print("  ")
		fmt.Print(media[i].Path)
		fmt.Print(": ")

		berr = !data.ChkHash(media[i].Path, media[i].Hash)

		if berr {
			errCount++
			fmt.Println("Error")
		} else {
			fmt.Println("OK")
		}
	}

	if errCount > 0 {
		fmt.Fprintln(os.Stderr, "Media integrity error detected!")
	}

	medhash := data.Medhash{
		Version: packageinfo.Version,
		Media:   media,
	}

	medhashJSON, err := json.MarshalIndent(medhash, "", "    ")

	if err != nil {
		panic(err)
	}

	medhashFile, err := os.Create("medhash.json")

	if err != nil {
		panic(err)
	}

	_, err = medhashFile.Write(medhashJSON)

	if err != nil {
		panic(err)
	}

	medhashFile.Sync()

	fmt.Println("Done!")
}
