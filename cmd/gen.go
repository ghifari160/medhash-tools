package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ghifari160/medhash-tools/color"
	"github.com/ghifari160/medhash-tools/medhash"
)

// Gen subcommand generates a MedHash manifest for the directories specified in Dirs.
type Gen struct {
	Dirs    []string `arg:"positional"`
	Ignores []string `arg:"--ignore,-i" help:"ignore patterns"`

	Ed25519Key string `arg:"--ed25519-privkey" help:"Path to Ed25519 private key"`

	CmdConfig
}

func (g *Gen) Execute() (status int) {
	var config medhash.Config

	if g.All {
		config = medhash.AllConfig
	} else if g.XXH3 || g.SHA512 || g.SHA3 || g.SHA256 || g.SHA1 || g.MD5 {
		config.XXH3 = g.XXH3
		config.SHA512 = g.SHA512
		config.SHA3 = g.SHA3
		config.SHA256 = g.SHA256
		config.SHA1 = g.SHA1
		config.MD5 = g.MD5
	} else if g.Default {
		config = medhash.DefaultConfig
	}

	if g.Ed25519 {
		if g.Ed25519Key == "" {
			confDir, err := ConfigDir()
			if err != nil {
				color.Printf("error: %v\n", err)
				status = 1

				return
			}

			g.Ed25519Key = filepath.Join(confDir, "ed25519.key")
		}

		keyStr, err := os.ReadFile(g.Ed25519Key)
		if err != nil {
			color.Printf("error: %v\n", err)
			status = 1

			return
		}

		privKey, private := medhash.DecodeKey(keyStr)
		if !private {
			color.Println("Expected private key. Got public key.")
			status = 1

			return
		}

		config.Ed25519.Enabled = true
		config.Ed25519.PrivateKey = privKey
	}

	if len(g.Dirs) < 1 {
		cwd, err := os.Getwd()
		if err != nil {
			color.Printf("error: %v\n", err)
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
		color.Printf("Generating MedHash for %s\n", dir)

		c := config
		c.Dir = dir

		errs := make([]error, 0)

		manifest, err := GenFunc(c, g.Ignores)
		errs = append(errs, err...)

		if g.Ed25519 {
			signed, err := SignFunc(c, manifest)
			errs = append(errs, err...)

			manifest = signed
		}

		errs = append(errs, WriteFunc(c, manifest)...)

		if len(errs) > 0 {
			color.Println(MsgFinalError)
			status = 1

			for _, err := range errs {
				color.Println(err)
			}
		} else {
			color.Println(MsgFinalDone)
		}
	}

	return
}

// GenFunc generates a Manifest using the provided config.
func GenFunc(config medhash.Config, ignores []string) (manifest *medhash.Manifest, errs []error) {
	media := make([]medhash.Media, 0)

	err := filepath.Walk(config.Dir, func(path string, info fs.FileInfo, err error) error {
		if !info.Mode().IsRegular() {
			return nil
		}

		color.Printf("  %s: ", path)

		c := config
		c.Dir = config.Dir

		if err != nil {
			color.Println(MsgStatusError)
			errs = append(errs, err)

			return nil
		}

		rel, err := filepath.Rel(c.Dir, path)
		if err != nil {
			color.Println(MsgStatusError)
			errs = append(errs, err)

			return nil
		}

		for _, ignore := range ignores {
			matched, err := filepath.Match(ignore, rel)
			if err != nil {
				color.Println(MsgStatusError)
				errs = append(errs, err)

				continue
			}

			if matched {
				color.Println(MsgStatusSkipped)

				return nil
			}
		}

		c.Path = rel

		med, err := medhash.GenHash(c)
		if err != nil {
			color.Println(MsgStatusError)
			errs = append(errs, err)

			return nil
		}

		color.Println(MsgStatusOK)
		media = append(media, med)

		return nil
	})
	if err != nil {
		errs = append(errs, err)
	}

	color.Println("Sanity checking files")

	for i, med := range media {
		c := config

		color.Printf("  %s: ", (filepath.Join(c.Dir, med.Path)))

		valid, err := medhash.ChkHash(c, med)
		if err == nil && valid {
			color.Println(MsgStatusOK)

			continue
		} else if err == nil && !valid {
			color.Println(MsgStatusError)
			errs = append(errs, fmt.Errorf("invalid hash for %s", med.Path))
		} else {
			color.Println(MsgStatusError)
			errs = append(errs, err)
		}

		m := media[:i]
		if i+1 < len(media) {
			m = append(m, media[i+1:]...)
		}
		media = m
	}

	if len(media) > 0 {
		manifest = medhash.NewWithConfig(config)
		manifest.Media = media
	}

	return
}

// SignFunc signs the manifest.
func SignFunc(config medhash.Config, manifest *medhash.Manifest) (signed *medhash.Manifest, errs []error) {
	color.Println("Signing manifest")

	signed, err := medhash.Sign(config, manifest)
	if err != nil {
		errs = append(errs, err)
	}

	return
}

// WriteFunc writes manifest into the Manifest file in config.Dir.
func WriteFunc(config medhash.Config, manifest *medhash.Manifest) (errs []error) {
	manFile, err := manifest.JSON()
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

	return
}
