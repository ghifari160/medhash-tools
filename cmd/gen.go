// MedHash Tools
// Copyright (c) 2023 GHIFARI160
// MIT License

package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ghifari160/medhash-tools/medhash"
)

// Gen subcommand generates a MedHash manifest for the directories specified in Dirs.
type Gen struct {
	Dirs    []string `arg:"positional,required" `
	Ignores []string `arg:"--ignore,-i"`

	Default bool `arg:"--default,-d" default:"true"`
	All     bool `arg:"--all,-a"`

	SHA3   bool `arg:"--sha3"`
	SHA256 bool `arg:"--sha256"`
	SHA1   bool `arg:"--sha1"`
	MD5    bool `arg:"--md5"`
}

func (g *Gen) Execute() (status int) {
	var config medhash.Config

	if g.Default {
		g.SHA3 = true
		g.SHA256 = true
		g.SHA1 = true
		g.MD5 = true
	} else if g.All {
		g.SHA3 = true
		g.SHA256 = true
		g.SHA1 = true
		g.MD5 = true
	}

	config.SHA3 = g.SHA3
	config.SHA256 = g.SHA256
	config.SHA1 = g.SHA1
	config.MD5 = g.MD5

	var manifestIgnored bool

	for _, ignore := range g.Ignores {
		if ignore == medhash.DefaultManifestName {
			manifestIgnored = true
		}
	}

	if !manifestIgnored {
		g.Ignores = append(g.Ignores, medhash.DefaultManifestName)
	}

	for _, dir := range g.Dirs {
		fmt.Printf("Generating MedHash for %s\n", dir)

		c := config
		c.Dir = dir

		errs := g.gen(c)

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

// gen generates a MedHash manifest using the provided config.
func (g *Gen) gen(config medhash.Config) (errs []error) {
	media := make([]medhash.Media, 0)

	err := filepath.Walk(config.Dir, func(path string, info fs.FileInfo, e error) error {
		if !info.Mode().IsRegular() {
			return nil
		}

		fmt.Printf("  %s: ", path)

		c := config
		c.Dir = config.Dir

		if e != nil {
			fmt.Println("ERROR")
			errs = append(errs, e)

			return nil
		}

		rel, err := filepath.Rel(c.Dir, path)
		if err != nil {
			fmt.Println("ERROR")
			errs = append(errs, err)

			return nil
		}

		for _, ignore := range g.Ignores {
			matched, err := filepath.Match(ignore, rel)
			if err != nil {
				fmt.Println("ERROR")
				errs = append(errs, err)

				continue
			}

			if matched {
				fmt.Println("SKIPPED")

				return nil
			}
		}

		c.Path = rel

		med, err := medhash.GenHash(c)
		if err != nil {
			fmt.Println("ERROR")
			errs = append(errs, err)

			return nil
		}

		fmt.Println("OK")
		media = append(media, med)

		return nil
	})
	if err != nil {
		errs = append(errs, err)
	}

	fmt.Println("Sanity checking files")

	for i, med := range media {
		fmt.Printf("  %s: ", med.Path)

		valid, err := medhash.ChkHash(config.Dir, med)
		if err == nil && valid {
			if err == nil {
				fmt.Println("OK")
				continue
			} else {
				fmt.Println("ERROR")
				errs = append(errs, err)
			}
		} else {
			fmt.Println("ERROR")
			errs = append(errs, err)
		}

		m := media[:i]
		if i+1 < len(media) {
			m = append(m, media[i+1:]...)
		}
		media = m
	}

	if len(media) > 0 {
		manifest := medhash.NewWithConfig(config)
		manifest.Media = media

		manFile, err := json.MarshalIndent(manifest, "", "  ")
		if err != nil {
			errs = append(errs, err)
			return
		}

		f, err := os.Create(filepath.Join(config.Dir, medhash.DefaultManifestName))
		if err != nil {
			errs = append(errs, err)
			return
		}
		defer f.Close()

		_, err = f.Write(manFile)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return
}
