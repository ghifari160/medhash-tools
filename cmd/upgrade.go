package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghifari160/medhash-tools/color"
	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/stretchr/objx"
)

// CurrentSpec is the most current Specification implemented.
const CurrentSpec = "0.5.0"

// Upgrade subcommand upgrades legacy Manifest to the current spec version.
type Upgrade struct {
	Dirs    []string `arg:"positional"`
	Ignores []string `arg:"--ignore,-i" help:"ignore patterns"`

	Force bool `arg:"--force" help:"force upgrade current Manifest"`

	CmdConfig
}

func (u *Upgrade) Execute() (status int) {
	var config medhash.Config

	if u.All {
		config = medhash.AllConfig
	} else if u.XXH3 || u.SHA3 || u.SHA256 || u.SHA1 || u.MD5 {
		config.XXH3 = u.XXH3
		config.SHA3 = u.SHA3
		config.SHA256 = u.SHA256
		config.SHA1 = u.SHA1
		config.MD5 = u.MD5
	} else if u.Default {
		config = medhash.DefaultConfig
	}

	if len(u.Dirs) < 1 {
		cwd, err := os.Getwd()
		if err != nil {
			color.Printf("error: %v\n", err)
			status = 1

			return
		}

		u.Dirs = append(u.Dirs, cwd)
	}

	var manifestIgnored bool

	for _, ignore := range u.Ignores {
		if ignore == medhash.DefaultManifestName {
			manifestIgnored = true
		}
	}

	if !manifestIgnored {
		u.Ignores = append(u.Ignores, medhash.DefaultManifestName)
	}

	for _, dir := range u.Dirs {
		color.Printf("Upgrading MedHash for %s\n", dir)

		c := config
		c.Dir = dir

		var errs []error

		_, err := os.Stat(filepath.Join(dir, medhash.DefaultManifestName))
		if os.IsNotExist(err) {
			_, err := os.Stat(filepath.Join(dir, "sums.txt"))
			if err != nil {
				errs = []error{err}
			} else {
				color.Println("Legacy Manifest detected!")
				errs = u.v010(c)
			}
		} else if err != nil {
			errs = []error{err}
		} else {
			errs = u.upgradeJSON(c)
		}

		if errs != nil {
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

// v010 upgrades a Manifest spec v0.1.0 to the current Manifest spec version.
func (u *Upgrade) v010(genConfig medhash.Config) (errs []error) {
	chkConfig := medhash.Config{
		Dir: genConfig.Dir,

		SHA256: true,
	}

	legacyPath := filepath.Join(genConfig.Dir, "sums.txt")

	legacy, err := os.ReadFile(legacyPath)
	if err != nil {
		errs = append(errs, err)
		return
	}
	u.Ignores = append(u.Ignores, "sums.txt")

	legacyMan := strings.Split(string(legacy), "\n")

	manifest := &medhash.Manifest{
		Media: make([]medhash.Media, 0),
	}
	manifest.Config = chkConfig

	color.Printf("Checking legacy Manifest for %s\n", genConfig.Dir)

	chkErrs := make([]error, 0)
	for i, med := range legacyMan {
		m := strings.Fields(med)

		if len(m) < 1 {
			continue
		} else if len(m) < 2 {
			color.Printf("Unknown media format (line %d): [%s]\n", i, strings.Join(m, ","))

			continue
		}

		newMed := medhash.Media{
			Path: m[0],
			Hash: medhash.Hash{
				SHA256: m[1],
			},
		}
		manifest.Media = append(manifest.Media, newMed)
	}

	for _, med := range manifest.Media {
		color.Printf("  %s: ", med.Path)
		err := manifest.Check(med.Path)
		if err != nil {
			chkErrs = append(chkErrs, err)
			color.Println(MsgStatusError)
		} else {
			color.Println(MsgStatusOK)
		}
	}
	errs = append(errs, chkErrs...)
	if len(chkErrs) > 0 {
		return
	}

	color.Printf("Generating MedHash for %s\n", genConfig.Dir)
	errs = append(errs, GenFunc(genConfig, u.Ignores)...)

	return
}

// upgradeJSON upgrades a Manifest to the current Manifest spec version.
func (u *Upgrade) upgradeJSON(config medhash.Config) (errs []error) {
	legacyPath := filepath.Join(config.Dir, medhash.DefaultManifestName)

	legacyFile, err := os.ReadFile(legacyPath)
	if err != nil {
		errs = append(errs, err)
		return
	}

	legacy, err := objx.FromJSON(string(legacyFile))
	if err != nil {
		errs = append(errs, err)
		return
	}

	ver := legacy.Get("version").Str()

	switch ver {
	case "0.2.0":
		color.Println("Manifest v0.2.0 detected!")
		errs = append(errs, u.upgradeV020(config, legacy)...)

	case "0.3.0":
		color.Println("Manifest v0.3.0 detected!")
		errs = append(errs, u.upgradeV030(config, legacy)...)

	case "0.4.0":
		color.Println("Manifest v0.4.0 detected!")
		errs = append(errs, u.upgradeV040(config, legacy)...)

	case "0.5.0":
		color.Println("Manifest v0.5.0 detected!")
		errs = append(errs, u.upgradeV050(config, legacy)...)

	default:
		errs = append(errs, fmt.Errorf("unexpected version: %v", legacy.Get("version").Data()))
		return
	}

	if len(errs) > 0 {
		return
	}

	color.Println("Generating MedHash")

	errs = append(errs, GenFunc(config, u.Ignores)...)

	return
}

// upgradeV020 upgrades a Manifest spec v0.2.0 to the current Manifest spec version.
func (u *Upgrade) upgradeV020(genConfig medhash.Config, legacy objx.Map) (errs []error) {
	chkConfig := medhash.Config{
		Dir: genConfig.Dir,

		SHA256: true,
	}

	color.Printf("Checking legacy manifest for %s\n", genConfig.Dir)

	if ver := legacy.Get("version").Str(); ver == "" || ver != "0.2.0" {
		errs = append(errs, fmt.Errorf("unexpected version: %v", legacy.Get("version").Data()))
		return
	}

	manifest, e := mapToManifest(legacy.Get("media"))
	errs = append(errs, e...)
	manifest.Config = chkConfig
	errs = append(errs, chkManifest(manifest)...)

	return
}

// upgradeV030 upgrades a Manifest spec v0.3.0 to the current Manifest spec version.
func (u *Upgrade) upgradeV030(genConfig medhash.Config, legacy objx.Map) (errs []error) {
	chkConfig := medhash.Config{
		Dir: genConfig.Dir,

		SHA256: true,
		SHA1:   true,
		MD5:    true,
	}

	color.Printf("Checking legacy manifest for %s\n", genConfig.Dir)

	if ver := legacy.Get("version").Str(); ver == "" || ver != "0.3.0" {
		errs = append(errs, fmt.Errorf("unexpected version: %v", legacy.Get("version").Data()))
		return
	}

	manifest, e := mapToManifest(legacy.Get("media"))
	errs = append(errs, e...)
	manifest.Config = chkConfig
	errs = append(errs, chkManifest(manifest)...)

	return
}

// upgradeV040 upgrades a Manifest spec v0.4.0 to the current Manifest spec version.
func (u *Upgrade) upgradeV040(genConfig medhash.Config, legacy objx.Map) (errs []error) {
	chkConfig := medhash.AllConfig
	chkConfig.Dir = genConfig.Dir

	if ver := legacy.Get("version").Str(); ver == "" || ver != "0.4.0" {
		errs = append(errs, fmt.Errorf("unexpected version: %v", legacy.Get("version").Data()))
		return
	}

	manifest, e := mapToManifest(legacy.Get("media"))
	errs = append(errs, e...)
	manifest.Config = chkConfig
	errs = append(errs, chkManifest(manifest)...)

	return
}

// upgradeV050 upgrades a Manifest spec v0.5.0 to the current Manifest spec version.
//
// With the Force flag enabled, this function regenerates a current Manifest with the config.
// Otherwise, this function is a placeholder.
// It does nothing.
func (u *Upgrade) upgradeV050(genConfig medhash.Config, legacy objx.Map) (errs []error) {
	chkConfig := medhash.AllConfig
	chkConfig.Dir = genConfig.Dir

	if ver := legacy.Get("version").Str(); ver == "" || ver != "0.5.0" {
		errs = append(errs, fmt.Errorf("unexpected version: %v", legacy.Get("version").Data()))
		return
	}

	if !u.Force {
		errs = append(errs, fmt.Errorf("manifest v0.5.0 is the current spec"))
		return
	}

	color.Printf("Forced to regenerate Manifest v0.5.0 for %s!\n", genConfig.Dir)

	manifest, e := mapToManifest(legacy.Get("media"))
	errs = append(errs, e...)
	manifest.Config = chkConfig
	errs = append(errs, chkManifest(manifest)...)

	return
}

// mapToManifest builds a medhash.Manifest from a JSON map.
func mapToManifest(legacyMed *objx.Value) (manifest *medhash.Manifest, errs []error) {
	if !legacyMed.IsInterSlice() {
		errs = append(errs, fmt.Errorf("invalid media array: %v", legacyMed.Data()))
		return
	}

	legacy := legacyMed.InterSlice()
	manifest = &medhash.Manifest{
		Media: make([]medhash.Media, len(legacy)),
	}

	for i, medInter := range legacy {
		msi, v := medInter.(map[string]any)
		if !v {
			errs = append(errs, fmt.Errorf("unknown media %d: %T %v", i, medInter, medInter))
			continue
		}
		med := objx.Map(msi)

		media := medhash.Media{}
		if path := med.Get("path").Str(); path == "" {
			errs = append(errs, fmt.Errorf("unexpected path for media %d: %v",
				i, med.Get("path").Data()))
			continue
		} else {
			media.Path = path
		}

		if !med.Get("hash.xxh3").IsStr() && !med.Get("hash.xxh3").IsNil() {
			errs = append(errs, hashErr{"xxh3", i, med.Get("hash.xxh3").Data()})
			continue
		} else {
			media.Hash.XXH3 = med.Get("hash.xxh3").Str()
		}

		if !med.Get("hash.sha512").IsStr() && !med.Get("hash.sha512").IsNil() {
			errs = append(errs, hashErr{"sha512", i, med.Get("hash.sha512").Data()})
			continue
		} else {
			media.Hash.SHA512 = med.Get("hash.sha512").Str()
		}

		if !med.Get("hash.sha3").IsStr() && !med.Get("hash.sha3").IsNil() {
			errs = append(errs, hashErr{"sha3", i, med.Get("hash.sha3").Data()})
			continue
		} else {
			media.Hash.SHA3 = med.Get("hash.sha3").Str()
		}
		if !med.Get("hash.sha3-256").IsStr() && !med.Get("hash.sha3-256").IsNil() {
			errs = append(errs, hashErr{"sha3-256", i, med.Get("hash.sha3-256").Data()})
			continue
		} else {
			media.Hash.SHA3_256 = med.Get("hash.sha3-256").Str()
		}

		if !med.Get("hash.sha256").IsStr() && !med.Get("hash.sha256").IsNil() {
			errs = append(errs, hashErr{"sha256", i, med.Get("hash.sha256").Data()})
			continue
		} else {
			media.Hash.SHA256 = med.Get("hash.sha256").Str()
		}

		if !med.Get("hash.sha1").IsStr() && !med.Get("hash.sha1").IsNil() {
			errs = append(errs, hashErr{"sha1", i, med.Get("hash.sha1").Data()})
			continue
		} else {
			media.Hash.SHA1 = med.Get("hash.sha1").Str()
		}

		if !med.Get("hash.md5").IsStr() && !med.Get("hash.md5").IsNil() {
			errs = append(errs, hashErr{"md5", i, med.Get("hash.md5").Data()})
			continue
		} else {
			media.Hash.MD5 = med.Get("hash.md5").Str()
		}

		manifest.Media[i] = media
	}
	return
}

// hashErr is an error type for invalid hash values.
type hashErr struct {
	alg   string
	index int
	data  any
}

func (err hashErr) Error() string {
	return fmt.Sprintf("unexpected %s for media %d: %v", err.alg, err.index, err.data)
}

// chkManifest verifies the Hashes for all Media in the provided manifest.
func chkManifest(manifest *medhash.Manifest) (errs []error) {
	for _, media := range manifest.Media {
		color.Printf("  %s: ", filepath.Join(manifest.Config.Dir, media.Path))
		err := manifest.Check(media.Path)
		if err != nil {
			errs = append(errs, err)
			color.Println(MsgStatusError)
		} else {
			color.Println(MsgStatusOK)
		}
	}
	return
}
