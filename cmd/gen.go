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
	Dirs    []string `arg:"positional"`
	Ignores []string `arg:"--ignore,-i" help:"ignore patterns"`

	Default bool `arg:"--default,-d" default:"true" help:"use default preset"`
	All     bool `arg:"--all,-a" help:"use all algorithms"`

	SHA3   bool `arg:"--sha3" help:"use SHA3-256"`
	SHA256 bool `arg:"--sha256" help:"use SHA256"`
	SHA1   bool `arg:"--sha1" help:"use SHA1"`
	MD5    bool `arg:"--md5" help:"use MD5"`
}

func (g *Gen) Execute() (status int) {
	var config medhash.Config

	if g.All {
		config = medhash.AllConfig
	} else if g.SHA3 || g.SHA256 || g.SHA1 || g.MD5 {
		config.SHA3 = g.SHA3
		config.SHA256 = g.SHA256
		config.SHA1 = g.SHA1
		config.MD5 = g.MD5
	} else if g.Default {
		config = medhash.DefaultConfig
	}

	if len(g.Dirs) < 1 {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			status = 1

			return
		}

		g.Dirs = append(g.Dirs, cwd)
	}

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

		errs := GenFunc(c, g.Ignores)

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

// GenFunc generates a Manifest using the provided config.
func GenFunc(config medhash.Config, ignores []string) (errs []error) {
	media := make([]medhash.Media, 0)

	err := filepath.Walk(config.Dir, func(path string, info fs.FileInfo, err error) error {
		if !info.Mode().IsRegular() {
			return nil
		}

		fmt.Printf("  %s: ", path)

		c := config
		c.Dir = config.Dir

		if err != nil {
			fmt.Println("ERROR")
			errs = append(errs, err)

			return nil
		}

		rel, err := filepath.Rel(c.Dir, path)
		if err != nil {
			fmt.Println("ERROR")
			errs = append(errs, err)

			return nil
		}

		for _, ignore := range ignores {
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
		fmt.Printf("  %s: ", (filepath.Join(config.Dir, med.Path)))

		valid, err := medhash.ChkHash(config.Dir, med)
		if err == nil && valid {
			fmt.Println("OK")

			continue
		} else if err == nil && !valid {
			fmt.Println("ERROR")
			errs = append(errs, fmt.Errorf("invalid hash for %s", med.Path))
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
