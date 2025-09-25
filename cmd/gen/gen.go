package gen

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ghifari160/medhash-tools/cmd"
	"github.com/ghifari160/medhash-tools/color"
	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/urfave/cli/v3"
)

func init() {
	cmd.RegisterCmd(CommandGen())
}

func CommandGen() *cli.Command {
	return &cli.Command{
		Name:  "gen",
		Usage: "generate MedHash Manifest",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "ignore",
				Aliases: []string{"i"},
				Usage:   "ignore patterns",
			},
		},
		MutuallyExclusiveFlags: []cli.MutuallyExclusiveFlags{
			{
				Flags: [][]cli.Flag{
					{
						&cli.BoolFlag{
							Name:  "default",
							Usage: "use default preset",
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
		Action: GenAction,
	}
}

func GenAction(ctx context.Context, command *cli.Command) error {
	config := cmd.ConfigFromFlags(command)

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
			color.Printf("[%d/%d] Generating MedHash for %s\n", i+1, len(dirs), dir)
		} else {
			color.Printf("Generating MedHash for %s\n", dir)
		}

		config := config
		config.Dir = dir

		err := GenFunc(config, ignores)
		if err != nil {
			errs = cmd.JoinErrors(errs, err)
		}
	}

	return cmd.FinalizeAction(errs)
}

// GenFunc generates a Manifest using the provided config.
func GenFunc(config medhash.Config, ignores []string) error {
	manifest, err := medhash.NewWithConfig(config)
	if err != nil {
		return err
	}

	var errs error
	err = filepath.Walk(config.Dir, func(path string, info fs.FileInfo, err error) error {
		if !info.Mode().IsRegular() {
			return nil
		}

		color.Printf("  %s: ", path)

		if err != nil {
			color.Println(cmd.MsgStatusError)
			errs = cmd.JoinErrors(errs, fmt.Errorf("cannot access %s: %w", path, err))
			return nil
		}

		c := config
		c.Dir = config.Dir

		rel, err := filepath.Rel(c.Dir, path)
		if err != nil {
			color.Println(cmd.MsgStatusError)
			errs = cmd.JoinErrors(errs, err)
		}

		for _, ignore := range ignores {
			matched, err := filepath.Match(ignore, rel)
			if err != nil {
				color.Println(cmd.MsgStatusError)
				errs = cmd.JoinErrors(errs, err)
			}

			if matched {
				color.Println(cmd.MsgStatusSkipped)
				return nil
			}
		}

		err = manifest.Add(rel)
		if err != nil {
			color.Println(cmd.MsgStatusError)
			errs = cmd.JoinErrors(errs, err)
		} else {
			color.Println(cmd.MsgStatusOK)
		}

		return nil
	})
	if err != nil {
		errs = cmd.JoinErrors(errs, err)
	}

	color.Println("Sanity checking files")

	for _, med := range manifest.Media {
		c := config

		color.Printf("  %s: ", (filepath.Join(c.Dir, med.Path)))

		err := manifest.Check(med.Path)
		if err != nil {
			errs = cmd.JoinErrors(errs, err)
			color.Println(cmd.MsgStatusError)
		} else {
			color.Println(cmd.MsgStatusOK)
		}
	}

	manFile, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		errs = cmd.JoinErrors(errs, err)
		return errs
	}

	f, err := os.Create(filepath.Join(config.Dir, medhash.DefaultManifestName))
	if err != nil {
		errs = cmd.JoinErrors(errs, err)
		return errs
	}
	defer f.Close()

	_, err = f.Write(manFile)
	if err != nil {
		errs = cmd.JoinErrors(errs, err)
	}
	return errs
}
