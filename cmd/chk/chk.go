package chk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"aead.dev/minisign"
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
			&cli.StringFlag{
				Name:  "minisign-keyfile",
				Usage: "verify the Manifest Minisign signature with this public key file",
			},
			&cli.StringFlag{
				Name:  "minisign-key",
				Usage: "verify the Manifest Minisign signature with this public key as Base64 string",
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

	if command.IsSet("minisign-keyfile") || command.IsSet("minisign-key") {
		var key minisign.PublicKey
		var err error
		if command.IsSet("minisign-keyfile") {
			key, err = minisign.PublicKeyFromFile(command.String("minisign-keyfile"))
		} else {
			err = key.UnmarshalText([]byte(command.String("minisign-key")))
		}
		if err != nil {
			return cmd.FinalizeAction(err)
		}
		config.Minisign.Enabled = true
		config.Minisign.PubKey = key
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

	var shouldVerifySig bool
	if config.Ed25519.Enabled {
		shouldVerifySig = true
	}
	if config.Minisign.Enabled {
		shouldVerifySig = true
	}

	if shouldVerifySig {
		color.Printf("Verifying manifest signature ")

		if config.Minisign.Enabled {
			sidecarFile := manPath + ".minisig"
			sig, err := os.ReadFile(sidecarFile)
			if err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					err = fmt.Errorf("cannot check for minisign sidecar file at %s", sidecarFile)
				} else {
					err = nil
				}
			} else {
				inManifest := strings.TrimSpace(manifest.Signature.Minisign)
				sig := strings.TrimSpace(string(sig))
				if inManifest == "" {
					manifest.Signature.Minisign = sig
				} else if inManifest != "" && inManifest != sig {
					err = fmt.Errorf("signature in manifest differs from sidecar file %s", sidecarFile)
				}
			}
			if err != nil {
				color.Println(cmd.MsgStatusError)
				return err
			}
		}

		err = manifest.Verify()
		if err != nil {
			color.Println(cmd.MsgStatusError)
			return err
		}

		color.Println(cmd.MsgStatusOK)
		color.Println("Checking files")
	}

	var errs error

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
