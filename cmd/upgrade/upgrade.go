package upgrade

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghifari160/medhash-tools/cmd"
	"github.com/ghifari160/medhash-tools/cmd/gen"
	"github.com/ghifari160/medhash-tools/color"
	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/stretchr/objx"
	"github.com/urfave/cli/v3"
)

// CurrentSpec is the most current Specification implemented.
const CurrentSpec = "0.5.0"

var (
	v010ChkConfig = medhash.Config{SHA256: true}
	v020ChkConfig = medhash.Config{SHA256: true}
	v030ChkConfig = medhash.Config{SHA256: true, SHA1: true, MD5: true}
	v040ChkConfig = medhash.AllConfig
	v050ChkConfig = medhash.AllConfig
)

func init() {
	cmd.RegisterCmd(CommandUpgrade())
}

func CommandUpgrade() *cli.Command {
	return &cli.Command{
		Name:  "upgrade",
		Usage: "upgrade MedHash Manifest",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "ignore",
				Aliases: []string{"i"},
				Usage:   "ignore patterns",
			},
			&cli.BoolFlag{
				Name:  "force",
				Usage: "force upgrade current Manifest",
			},
		},
		MutuallyExclusiveFlags: []cli.MutuallyExclusiveFlags{
			{
				Flags: [][]cli.Flag{
					{
						&cli.BoolFlag{
							Name:  "default",
							Usage: "use default preset",
							Value: true,
						},
					},
					{
						&cli.BoolFlag{
							Name:  "all",
							Usage: "use all algorithms",
						},
					},
					cmd.HashAlgs(),
				},
			},
		},
		Action: UpgradeAction,
	}
}

func UpgradeAction(ctx context.Context, command *cli.Command) error {
	var config medhash.Config

	if command.Bool("all") {
		config = medhash.AllConfig
	} else if command.Bool("default") {
		config = medhash.DefaultConfig
	} else {
		config.XXH3 = command.Bool("xxh3")
		config.SHA512 = command.Bool("sha512")
		config.SHA3 = command.Bool("sha3")
		config.SHA256 = command.Bool("sha256")
		config.SHA1 = command.Bool("sha1")
		config.MD5 = command.Bool("md5")
	}

	force := command.Bool("force")

	dirs := command.Args().Slice()
	if len(dirs) < 1 {
		cwd, err := os.Getwd()
		if err != nil {
			return cli.Exit(fmt.Errorf("cannot get working directory: %w", err), 1)
		}
		dirs = append(dirs, cwd)
	}

	var manifestIgnored bool
	ignores := command.StringSlice("ignore")

	for _, ignore := range ignores {
		if ignore == medhash.DefaultManifestName {
			manifestIgnored = true
		}
	}

	if !manifestIgnored {
		ignores = append(ignores, medhash.DefaultManifestName)
	}

	var errs error
	for i, dir := range dirs {
		if len(dirs) > 1 {
			color.Printf("[%d/%d] Upgrading MedHash for %s\n", i+1, len(dirs), dir)
		} else {
			color.Printf("Upgrading MedHash for %s\n", dir)
		}

		conf := config
		conf.Dir = dir

		_, err := os.Stat(filepath.Join(dir, medhash.DefaultManifestName))
		if errors.Is(err, os.ErrNotExist) {
			_, err := os.Stat(filepath.Join(dir, "sums.txt"))
			if errors.Is(err, os.ErrNotExist) {
				errs = cmd.JoinErrors(errs, fmt.Errorf("no manifest.json or sums.txt found in %s", dir))
			} else if err != nil {
				errs = cmd.JoinErrors(errs, err)
			} else {
				color.Println("Legacy Manifest detected!")
				errs = cmd.JoinErrors(errs, upgradeV010(conf, ignores, force))
			}
		} else if err != nil {
			errs = cmd.JoinErrors(errs, err)
		} else {
			errs = cmd.JoinErrors(errs, upgradeJSON(conf, ignores, force))
		}
	}

	if errs != nil {
		color.Println(cmd.MsgFinalError)
		for _, err := range cmd.UnwrapJoinedErrors(errs) {
			color.Println(err)
		}
		return cli.Exit("", 1)
	}

	color.Println(cmd.MsgFinalDone)
	return nil
}

// upgradeV010 upgrades a Manifest spec v0.1.0 to the current Manifest spec version.
func upgradeV010(genConfig medhash.Config, ignores []string, force bool) error {
	chkConfig := v010ChkConfig
	chkConfig.Dir = genConfig.Dir

	legacyPath := filepath.Join(chkConfig.Dir, "sums.txt")

	var errs error

	legacyManifest, err := os.ReadFile(legacyPath)
	if err != nil {
		return err
	}
	ignores = append(ignores, "sums.txt")

	legacyMedias := strings.Split(string(legacyManifest), "\n")

	convertedManifest := &medhash.Manifest{
		Media: make([]medhash.Media, 0),
	}
	convertedManifest.Config = chkConfig

	color.Printf("Checking legacy Manifest for %s\n", convertedManifest.Config.Dir)

	var chkErrs error
	for i, med := range legacyMedias {
		legacyMedia := strings.Fields(med)
		if len(legacyMedia) < 1 {
			continue
		} else if len(legacyMedia) < 2 {
			color.Printf("Unknown media format (line %d): [%s]\n", i, strings.Join(legacyMedia, ","))
			continue
		}

		convertedMedia := medhash.Media{
			Path: legacyMedia[0],
			Hash: medhash.Hash{
				SHA256: legacyMedia[1],
			},
		}
		convertedManifest.Media = append(convertedManifest.Media, convertedMedia)
	}

	for _, media := range convertedManifest.Media {
		color.Printf("  %s: ", media.Path)
		err := media.Check(convertedManifest.Config)
		if err != nil {
			chkErrs = cmd.JoinErrors(chkErrs, err)
			color.Print(cmd.MsgStatusError)
		} else {
			color.Println(cmd.MsgStatusOK)
		}
	}
	if chkErrs != nil {
		errs = cmd.JoinErrors(errs, chkErrs)
		return errs
	}

	color.Printf("Generating MedHash for %s\n", genConfig.Dir)
	errs = cmd.JoinErrors(errs, gen.GenFunc(genConfig, ignores))
	return errs
}

func upgradeJSON(genConfig medhash.Config, ignores []string, force bool) error {
	var errs error
	legacyPath := filepath.Join(genConfig.Dir, medhash.DefaultManifestName)

	legacyFile, err := os.ReadFile(legacyPath)
	if err != nil {
		return err
	}

	legacyManifest, err := objx.FromJSON(string(legacyFile))
	if err != nil {
		return err
	}

	version := legacyManifest.Get("version").Str()
	switch version {
	case "0.2.0":
		color.Println("Manifest v0.2.0 detected!")
		errs = cmd.JoinErrors(errs, upgradeV020(genConfig, legacyManifest, force))

	case "0.3.0":
		color.Println("Manifest v0.3.0 detected!")
		errs = cmd.JoinErrors(errs, upgradeV030(genConfig, legacyManifest, force))

	case "0.4.0":
		color.Println("Manifest v0.4.0 detected!")
		errs = cmd.JoinErrors(errs, upgradeV040(genConfig, legacyManifest, force))

	case "0.5.0":
		color.Println("Manifest v0.5.0 detected!")
		errs = cmd.JoinErrors(errs, upgradeV050(genConfig, legacyManifest, force))

	default:
		errs = cmd.JoinErrors(errs, fmt.Errorf("unexpected version: %v", legacyManifest.Get("version").Data()))
	}
	if errs != nil {
		return errs
	}

	color.Println("Generating MedHash")

	return cmd.JoinErrors(errs, gen.GenFunc(genConfig, ignores))
}

// upgradeV020 upgrades a Manifest spec v0.2.0 to the current Manifest spec version.
func upgradeV020(genConfig medhash.Config, legacy objx.Map, force bool) error {
	var errs error
	chkConfig := v020ChkConfig
	chkConfig.Dir = genConfig.Dir

	color.Printf("Checking legacy manifest for %s\n", chkConfig.Dir)

	if err := expectVersion("0.2.0", legacy.Get("version")); err != nil {
		return err
	}

	convertedManifest, mapErrs := mapToManifest(legacy.Get("media"))
	errs = cmd.JoinErrors(errs, mapErrs)
	if errs != nil {
		return errs
	}
	convertedManifest.Config = chkConfig
	errs = cmd.JoinErrors(errs, chkManifest(convertedManifest))

	return errs
}

// upgradeV030 upgrades a Manifest spec v0.3.0 to the current Manifest spec version.
func upgradeV030(genConfig medhash.Config, legacy objx.Map, force bool) error {
	var errs error
	chkConfig := v030ChkConfig
	chkConfig.Dir = genConfig.Dir

	color.Printf("Checking legacy manifest for %s\n", chkConfig.Dir)

	if err := expectVersion("0.3.0", legacy.Get("version")); err != nil {
		return err
	}

	convertedManifest, mapErrs := mapToManifest(legacy.Get("media"))
	errs = cmd.JoinErrors(errs, mapErrs)
	if errs != nil {
		return errs
	}
	convertedManifest.Config = chkConfig
	errs = cmd.JoinErrors(errs, chkManifest(convertedManifest))

	return errs
}

// upgradeV040 upgrades a Manifest spec v0.4.0 to the current Manifest spec version.
func upgradeV040(genConfig medhash.Config, legacy objx.Map, force bool) error {
	var errs error
	chkConfig := v040ChkConfig
	chkConfig.Dir = genConfig.Dir

	color.Printf("Checking legacy manifest for %s\n", chkConfig.Dir)

	if err := expectVersion("0.4.0", legacy.Get("version")); err != nil {
		return err
	}

	convertedManifest, mapErrs := mapToManifest(legacy.Get("media"))
	errs = cmd.JoinErrors(errs, mapErrs)
	if errs != nil {
		return errs
	}
	convertedManifest.Config = chkConfig
	errs = cmd.JoinErrors(errs, chkManifest(convertedManifest))

	return errs
}

// upgradeV050 upgrades a Manifest spec v0.5.0 to the current Manifest spec version.
//
// With the force flag enabled, this function regenerates a current Manifest with the config.
// Otherwise, this function is a placeholder.
// It does nothing.
func upgradeV050(genConfig medhash.Config, legacy objx.Map, force bool) error {
	var errs error
	chkConfig := v050ChkConfig
	chkConfig.Dir = genConfig.Dir

	if err := expectVersion("0.5.0", legacy.Get("version")); err != nil {
		return err
	}

	if !force {
		return fmt.Errorf("manifest v0.5.0 is the current spec")
	}

	color.Printf("Forced to regenerated Manifest v0.5.0 %s!\n", chkConfig.Dir)

	convertedManifest, mapErrs := mapToManifest(legacy.Get("media"))
	errs = cmd.JoinErrors(errs, mapErrs)
	if errs != nil {
		return errs
	}
	convertedManifest.Config = chkConfig
	errs = cmd.JoinErrors(errs, chkManifest(convertedManifest))

	return errs
}

func expectVersion(expected string, actual *objx.Value) error {
	if ver := actual.Str(); ver != expected {
		return fmt.Errorf("unexpected version: %v", actual.Data())
	} else {
		return nil
	}
}

func mapToManifest(legacyMedias *objx.Value) (convertedManifest *medhash.Manifest, errs error) {
	if !legacyMedias.IsInterSlice() {
		errs = fmt.Errorf("invalid media array: %v", legacyMedias.Data())
		return
	}

	legacyMeds := legacyMedias.InterSlice()
	convertedManifest = &medhash.Manifest{
		Media: make([]medhash.Media, len(legacyMeds)),
	}

	for i, medInter := range legacyMeds {
		msi, v := medInter.(map[string]any)
		if !v {
			errs = cmd.JoinErrors(errs, fmt.Errorf("unknown media %d: %T %v", i, medInter, medInter))
			continue
		}
		media := objx.Map(msi)

		convertedMedia := medhash.Media{}
		if path := media.Get("path").Str(); path == "" {
			errs = cmd.JoinErrors(errs, fmt.Errorf("unexpected path for media %d: %v",
				i, media.Get("path").Data()))
			continue
		} else {
			convertedMedia.Path = path
		}

		hashErrs := cmd.JoinErrors(
			convertStrField(i, "xxh3", media.Get("hash.xxh3"), &convertedMedia.Hash.XXH3),
			convertStrField(i, "sha512", media.Get("hash.sha512"), &convertedMedia.Hash.SHA512),
			convertStrField(i, "sha3", media.Get("hash.sha3"), &convertedMedia.Hash.SHA3),
			convertStrField(i, "sha3_256", media.Get("hash.sha3_256"), &convertedMedia.Hash.SHA3_256),
			convertStrField(i, "sha256", media.Get("hash.sha256"), &convertedMedia.Hash.SHA256),
			convertStrField(i, "sha1", media.Get("hash.sha1"), &convertedMedia.Hash.SHA1),
			convertStrField(i, "md5", media.Get("hash.md5"), &convertedMedia.Hash.MD5),
		)
		if hashErrs != nil {
			errs = cmd.JoinErrors(errs, hashErrs)
			continue
		}
		convertedManifest.Media[i] = convertedMedia
	}

	return
}

// chkManifest verifies the Hashes for all Media in the provided manifest.
func chkManifest(manifest *medhash.Manifest) (errs error) {
	for _, media := range manifest.Media {
		color.Printf("  %s: ", filepath.Join(manifest.Config.Dir, media.Path))
		err := media.Check(manifest.Config)
		if err != nil {
			errs = cmd.JoinErrors(errs, err)
			color.Println(cmd.MsgStatusError)
		} else {
			color.Println(cmd.MsgStatusOK)
		}
	}
	return
}

func convertStrField(i int, field string, source *objx.Value, target *string) error {
	if !source.IsStr() && !source.IsNil() {
		return hashErr{field, i, source.Data()}
	}
	*target = source.Str()
	return nil
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
