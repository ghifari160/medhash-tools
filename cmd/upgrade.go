package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/stretchr/objx"
)

// Upgrade subcommand upgrades legacy Manifest to the current spec version.
type Upgrade struct {
	Dirs    []string `arg:"positional"`
	Ignores []string `arg:"--ignore,-i" help:"ignore patterns"`

	Force   bool `arg:"--force" help:"force upgrade current Manifest"`
	Default bool `arg:"--default,-d" default:"true" help:"use default preset"`
	All     bool `arg:"--all,-a" help:"use all algorithms"`

	SHA3   bool `arg:"--sha3" help:"use SHA3-256"`
	SHA256 bool `arg:"--sha256" help:"use SHA256"`
	SHA1   bool `arg:"--sha1" help:"use SHA1"`
	MD5    bool `arg:"--md5" help:"use MD5"`
}

func (u *Upgrade) Execute() (status int) {
	var config medhash.Config

	if u.All {
		config = medhash.AllConfig
	} else if u.SHA3 || u.SHA256 || u.SHA1 || u.MD5 {
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
			fmt.Printf("error: %v\n", err)
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
		fmt.Printf("Upgrading MedHash for %s\n", dir)

		c := config
		c.Dir = dir

		var errs []error

		_, err := os.Stat(filepath.Join(dir, medhash.DefaultManifestName))
		if os.IsNotExist(err) {
			_, err := os.Stat(filepath.Join(dir, "sums.txt"))
			if err != nil {
				errs = []error{err}
			} else {
				fmt.Println("Legacy Manifest detected!")
				errs = u.v010(c)
			}
		} else if err != nil {
			errs = []error{err}
		} else {
			errs = u.upgradeJSON(c)
		}

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

// v010 upgrades a Manifest spec v0.1.0 to the current Manifest spec version.
func (u *Upgrade) v010(config medhash.Config) (errs []error) {
	legacyPath := filepath.Join(config.Dir, "sums.txt")

	legacy, err := os.ReadFile(legacyPath)
	if err != nil {
		errs = append(errs, err)
		return
	}

	legacyMan := strings.Split(string(legacy), "\n")

	fmt.Printf("Checking legacy Manifest for %s\n", config.Dir)

	chkErrs := make([]error, 0)
	for i, med := range legacyMan {
		m := strings.Fields(med)

		if len(m) < 1 {
			continue
		} else if len(m) < 2 {
			fmt.Printf("Unknown media format (line %d): [%s]\n", i, strings.Join(m, ","))

			continue
		}

		newMed := medhash.Media{
			Path: m[0],
			Hash: medhash.Hash{
				SHA256: m[1],
			},
		}

		err = chkMedia(config.Dir, newMed)
		if err != nil {
			errs = append(errs, err)

			continue
		}
	}

	errs = append(errs, chkErrs...)
	if len(chkErrs) > 0 {
		return
	}

	fmt.Printf("Generating MedHash for %s\n", config.Dir)
	errs = append(errs, GenFunc(config, u.Ignores)...)

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

	switch breakoutSemver(ver) {
	case [3]int{0, 2, 0}:
		fmt.Println("Manifest v0.2.0 detected!")
		errs = append(errs, u.upgradeV020(config, legacy)...)

	case [3]int{0, 3, 0}:
		fmt.Println("Manifest v0.3.0 detected!")
		errs = append(errs, u.upgradeV030(config, legacy)...)

	case [3]int{0, 4, 0}:
		fmt.Println("Manifest v0.4.0 detected!")
		errs = append(errs, u.upgradeV040(config, legacy)...)

	default:
		errs = append(errs, fmt.Errorf("unexpected version: %v", legacy.Get("version").Data()))
		return
	}

	if len(errs) > 0 {
		return
	}

	fmt.Println("Generating MedHash")

	errs = append(errs, GenFunc(config, u.Ignores)...)

	return
}

// upgradeV020 upgrades a Manifest spec v0.2.0 to the current Manifest spec version.
func (u *Upgrade) upgradeV020(config medhash.Config, legacy objx.Map) (errs []error) {
	fmt.Printf("Checking legacy manifest for %s\n", config.Dir)

	if ver := legacy.Get("version").Str(); ver == "" || breakoutSemver(ver) != [3]int{0, 2, 0} {
		errs = append(errs, fmt.Errorf("unexpected version: %v", legacy.Get("version").Data()))
		return
	}

	media, e := mapToMedia(legacy.Get("media"))
	errs = append(errs, e...)

	errs = append(errs, chkMediaSlice(config.Dir, media)...)

	return
}

// upgradeV030 upgrades a Manifest spec v0.3.0 to the current Manifest spec version.
func (u *Upgrade) upgradeV030(config medhash.Config, legacy objx.Map) (errs []error) {
	fmt.Printf("Checking legacy manifest for %s\n", config.Dir)

	if ver := legacy.Get("version").Str(); ver == "" || breakoutSemver(ver) != [3]int{0, 3, 0} {
		errs = append(errs, fmt.Errorf("unexpected version: %v", legacy.Get("version").Data()))
		return
	}

	media, e := mapToMedia(legacy.Get("media"))
	errs = append(errs, e...)

	errs = append(errs, chkMediaSlice(config.Dir, media)...)

	return
}

// upgradeV040 upgrades a Manifest spec v0.4.0 to the current Manifest spec version.
//
// With the Force flag enabled, this function regenerates a current Manifest with the config.
// Otherwise, this function is a placeholder.
// It does nothing.
func (u *Upgrade) upgradeV040(config medhash.Config, legacy objx.Map) (errs []error) {
	if ver := legacy.Get("version").Str(); ver == "" || breakoutSemver(ver) != [3]int{0, 4, 0} {
		errs = append(errs, fmt.Errorf("unexpected version: %v", legacy.Get("version").Data()))
		return
	}

	if !u.Force {
		errs = append(errs, fmt.Errorf("manifest v0.4.0 is the current spec"))
		return
	}

	fmt.Printf("Forced to regenerate Manifest v0.4.0 for %s!\n", config.Dir)

	media, e := mapToMedia(legacy.Get("media"))
	errs = append(errs, e...)

	errs = append(errs, chkMediaSlice(config.Dir, media)...)

	return
}

// mapToMedia converts an objx.Value to medhash.Media slice.
func mapToMedia(legacyMed *objx.Value) (media []medhash.Media, errs []error) {
	if !legacyMed.IsInterSlice() {
		errs = append(errs, fmt.Errorf("invalid media array: %v", legacyMed.Data()))
		return
	}

	legacy := legacyMed.InterSlice()
	media = make([]medhash.Media, len(legacy))

	for i, medInter := range legacy {
		msi, v := medInter.(map[string]interface{})
		if !v {
			errs = append(errs, fmt.Errorf("unknown media %d: %T %v", i, medInter, medInter))
			continue
		}

		med := objx.Map(msi)

		m := medhash.Media{}

		if path := med.Get("path").Str(); path == "" {
			errs = append(errs, fmt.Errorf("unexpected path for media %d: %v",
				i, med.Get("path").Data()))

			continue
		} else {
			m.Path = path
		}

		if !med.Get("hash.sha3-256").IsStr() && !med.Get("hash.sha3-256").IsNil() {
			errs = append(errs, fmt.Errorf("unexpected sha3 for media %d: %v",
				i, med.Get("hash.sha3-256").Data()))

			continue
		} else {
			m.Hash.SHA3_256 = med.Get("hash.sha3-256").Str()
		}

		if !med.Get("hash.sha256").IsStr() && !med.Get("hash.sha256").IsNil() {
			errs = append(errs, fmt.Errorf("unexpected sha256 for media %d: %v",
				i, med.Get("hash.sha256").Data()))

			continue
		} else {
			m.Hash.SHA256 = med.Get("hash.sha256").Str()
		}

		if !med.Get("hash.sha1").IsStr() && !med.Get("hash.sha1").IsNil() {
			errs = append(errs, fmt.Errorf("unexpected sha1 for media %d: %v",
				i, med.Get("hash.sha1").Data()))

			continue
		} else {
			m.Hash.SHA1 = med.Get("hash.sha1").Str()
		}

		if !med.Get("hash.md5").IsStr() && !med.Get("hash.md5").IsNil() {
			errs = append(errs, fmt.Errorf("unexpected md5 for media %d: %v",
				i, med.Get("hash.md5").Data()))

			continue
		} else {
			m.Hash.MD5 = med.Get("hash.md5").Str()
		}

		media[i] = m
	}

	return
}

// chkMediaSlice verifies the Hashes for all Media in the provided slice.
func chkMediaSlice(dir string, media []medhash.Media) (errs []error) {
	for _, med := range media {
		err := chkMedia(dir, med)
		if err != nil {
			errs = append(errs, err)

			continue
		}
	}

	return
}

// chkMedia verifies Hashes for the Media.
func chkMedia(dir string, media medhash.Media) (err error) {
	fmt.Printf("  %s: ", filepath.Join(dir, media.Path))

	valid, err := medhash.ChkHash(dir, media)
	if err != nil {
		fmt.Println("ERROR")

		return
	} else if !valid {
		fmt.Println("ERROR")
		err = fmt.Errorf("invalid hash for %s", filepath.Join(dir, media.Path))

		return
	}

	fmt.Println("OK")

	return
}

// compareSemver compares two SemVer strings.
// It returns 1 if a > b, -1 if a < b, and 0 if a == b.
func compareSemver(a, b string) int {
	aInt := breakoutSemver(a)
	bInt := breakoutSemver(b)

	if aInt[0] > bInt[0] {
		return 1
	} else if aInt[0] < bInt[0] {
		return -1
	} else {
		if aInt[1] > bInt[1] {
			return 1
		} else if aInt[1] < bInt[1] {
			return -1
		} else {
			if aInt[2] > bInt[2] {
				return 1
			} else if aInt[2] < bInt[2] {
				return -1
			} else {
				return 0
			}
		}
	}
}

// breakoutSemver parses a SemVer string into an int array.
func breakoutSemver(semver string) (ver [3]int) {
	pattern := regexp.MustCompile(`(?:v{0,1})([0-9]*)(?:\.{0,1})([0-9]*)(?:\.{0,1})([0-9]*)`)

	verStr := make([]string, 3)

	verStr[0] = pattern.ReplaceAllString(semver, "$1")
	verStr[1] = pattern.ReplaceAllString(semver, "$2")
	verStr[2] = pattern.ReplaceAllString(semver, "$3")

	for i, str := range verStr {
		var err error

		ver[i], err = strconv.Atoi(str)
		if err != nil {
			ver[i] = 0
		}
	}

	return
}
