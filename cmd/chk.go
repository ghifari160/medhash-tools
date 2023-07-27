package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghifari160/medhash-tools/medhash"
)

// Chk subcommand verifies the directories specified in Dirs and for the files specified in Files.
type Chk struct {
	Dirs  []string `arg:"positional"`
	Files []string `arg:"--file,-f,separate" help:"check only these files"`

	Manifest string `arg:"--manifest,-m" help:"use this manifest"`

	CmdConfig
}

func (c *Chk) Execute() (status int) {
	var config medhash.Config

	if c.All {
		config = medhash.AllConfig
	} else if c.SHA3 || c.SHA256 || c.SHA1 || c.MD5 {
		config.SHA3 = c.SHA3
		config.SHA256 = c.SHA256
		config.SHA1 = c.SHA1
		config.MD5 = c.MD5
	} else if c.Default {
		config = medhash.DefaultConfig
	}

	if len(c.Dirs) < 1 {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			status = 1

			return
		}

		c.Dirs = append(c.Dirs, cwd)
	}

	for _, dir := range c.Dirs {
		conf := config
		conf.Dir = dir

		var manPath string

		if len(c.Manifest) > 0 {
			manPath = c.Manifest
		} else {
			manPath = filepath.Join(dir, medhash.DefaultManifestName)
		}

		fmt.Printf("Checking MedHash for %s\n", dir)

		errs := c.chk(manPath, conf, c.Files)
		if errs != nil {
			fmt.Println("Error!")
			status = 1

			for _, err := range errs {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Done!")
		}
	}

	return
}

func (c *Chk) chk(manPath string, config medhash.Config, files []string) (errs []error) {
	manFile, err := os.ReadFile(manPath)
	if err != nil {
		errs = append(errs, err)
		return
	}

	var manifest medhash.Manifest

	err = json.Unmarshal(manFile, &manifest)
	if err != nil {
		errs = append(errs, err)
		return
	}

	for _, med := range manifest.Media {
		fmt.Printf("  %s: ", filepath.Join(config.Dir, med.Path))

		if len(files) > 0 {
			skipped := true

			for _, file := range files {
				matched, err := filepath.Match(file, med.Path)
				if err != nil {
					fmt.Println("ERROR")
					errs = append(errs, err)

					continue
				}

				if matched {
					skipped = false
					break
				}
			}

			if skipped {
				fmt.Println("SKIPPED")
				continue
			}
		}

		valid, err := medhash.ChkHash(config, med)
		if err == nil && valid {
			fmt.Println("OK")
			continue
		} else {
			fmt.Println("ERROR")

			if err == nil {
				err = fmt.Errorf("invalid hash for %s", filepath.Join(config.Dir, med.Path))
			}

			errs = append(errs, err)
		}
	}

	return
}
