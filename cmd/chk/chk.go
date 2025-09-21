package chk

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghifari160/medhash-tools/cmd"
	"github.com/ghifari160/medhash-tools/color"
	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/urfave/cli/v3"
)

func init() {
	cmd.RegisterCmd(CommandChk())
}

func CommandChk() *cli.Command {
	return &cli.Command{
		Name:  "chk",
		Usage: "verify directories or files",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "check only these files",
			},
			&cli.StringFlag{
				Name:    "manifest",
				Aliases: []string{"m"},
				Usage:   "use this manifest",
			},
			&cli.StringFlag{
				Name:  "ed25519",
				Usage: "verify the Manifest Ed25519 signature with this public key",
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
		Action: ChkAction,
	}
}

func ChkAction(ctx context.Context, command *cli.Command) error {
	config := cmd.ConfigFromFlags(command)

	if command.IsSet("ed25519") {
		key, err := cmd.Loader(command.String("ed25519"), "ED25519 PUBLIC KEY")
		if err != nil {
			return cmd.FinalizeAction(err)
		}
		config.Ed25519.Enabled = true
		config.Ed25519.PubKey = key
	}

	dirs := command.Args().Slice()
	if len(dirs) < 1 {
		cwd, err := os.Getwd()
		if err != nil {
			return cli.Exit(fmt.Errorf("cannot get working directory: %w", err), 1)
		}
		dirs = append(dirs, cwd)
	}

	var errs error
	for i, dir := range dirs {
		conf := config
		conf.Dir = dir

		manPath := command.String("manifest")
		if manPath == "" {
			manPath = filepath.Join(dir, medhash.DefaultManifestName)
		}

		if len(dirs) > 1 {
			color.Printf("[%d/%d] Checking MedHash for %s\n", i+1, len(dirs), dir)
		} else {
			color.Printf("Checking MedHash for %s\n", dir)
		}

		errs = cmd.JoinErrors(errs, chk(manPath, conf, command.StringSlice("files")))
	}

	return cmd.FinalizeAction(errs)
}

func chk(manPath string, config medhash.Config, files []string) error {
	manFile, err := os.ReadFile(manPath)
	if err != nil {
		return err
	}

	var manifest medhash.Manifest

	err = json.Unmarshal(manFile, &manifest)
	if err != nil {
		return err
	}
	manifest.Config = config

	var errs error

	if config.Ed25519.Enabled {
		color.Printf("Verifying manifest signature ")
		err = manifest.Verify()
		if err != nil {
			color.Println(cmd.MsgStatusError)
			return err
		} else {
			color.Println(cmd.MsgStatusOK)
		}
		color.Println("Checking files")
	}

	for _, med := range manifest.Media {
		color.Printf("  %s: ", filepath.Join(config.Dir, med.Path))

		if len(files) > 0 {
			skipped := true

			for _, file := range files {
				matched, err := filepath.Match(file, med.Path)
				if err != nil {
					color.Println(cmd.MsgStatusError)
					errs = cmd.JoinErrors(errs, err)

					continue
				}

				if matched {
					skipped = false
					break
				}
			}

			if skipped {
				color.Println(cmd.MsgStatusSkipped)
				continue
			}
		}

		err := med.Check(manifest.Config)
		if err != nil {
			errs = cmd.JoinErrors(errs, err)
			color.Println(cmd.MsgStatusError)
		} else {
			color.Println(cmd.MsgStatusOK)
		}
	}

	return errs
}
