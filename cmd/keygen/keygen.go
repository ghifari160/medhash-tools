package keygen

import (
	"context"
	"crypto/ed25519"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"

	"aead.dev/minisign"
	"github.com/ghifari160/medhash-tools/cmd"
	"github.com/ghifari160/medhash-tools/color"
	"github.com/urfave/cli/v3"
)

func init() {
	RegisterSubcommand(
		Subcommand("ed25519", "generate Ed25519 keypair for signing MedHash Manifest",
			func(_ string) (pubKey, privKey []byte, err error) {
				return ed25519.GenerateKey(nil)
			},
			func(private bool, path string, data []byte) error {
				block := pem.Block{
					Bytes: data,
				}
				if private {
					block.Type = "ED25519 PRIVATE KEY"
				} else {
					block.Type = "ED25519 PUBLIC KEY"
				}

				f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
				if err != nil {
					return err
				}
				defer f.Close()

				return pem.Encode(f, &block)
			}))
	RegisterSubcommand(
		Subcommand("minisign", "generate Minisign keypair for signing MedHash Manifest",
			func(password string) (pubKeyData, privKeyData []byte, err error) {
				pubKey, privKey, err := minisign.GenerateKey(nil)
				if err != nil {
					return
				}

				if password != "" {
					privKeyData, err = minisign.EncryptKey(password, privKey)
				} else {
					privKeyData, err = privKey.MarshalText()
				}
				if err != nil {
					return
				}

				pubKeyData, err = pubKey.MarshalText()
				if err != nil {
					return
				}

				return
			},
			func(private bool, path string, data []byte) error {
				return os.WriteFile(path, data, 0600)
			}),
	)
	cmd.RegisterCmd(Command())
}

var subcommands []*cli.Command

// RegisterSubcommand registers a subcommand.
func RegisterSubcommand(command *cli.Command) {
	subcommands = append(subcommands, command)
}

// Subcommands returns all registered subcommands.
func Subcommands() []*cli.Command {
	return subcommands
}

// Command returns the Keygen command.
// Command must be called *after* all subcommands are registered.
func Command() *cli.Command {
	return &cli.Command{
		Name:     "keygen",
		Usage:    "generate keypair for signing MedHash Manifest",
		Commands: Subcommands(),
	}
}

// Subcommand creates a new subcommand for a specific keypair algorithm.
// It is up to the caller to register the returned subcommand.
func Subcommand(name, usage string, generator cmd.Generator, storer cmd.Storer) *cli.Command {
	return &cli.Command{
		Name:  name,
		Usage: usage,
		Flags: CommonFlags(),
		Action: func(ctx context.Context, command *cli.Command) error {
			return subcommandAction(ctx, command, generator, storer)
		},
	}
}

// CommonFlags returns a set of flags common to all subcommands.
func CommonFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "private",
			Usage: "path to store private key",
		},
		&cli.StringFlag{
			Name:  "public",
			Usage: "path to store public key",
		},
		&cli.BoolFlag{
			Name:  "force",
			Usage: "force generation and overwrite key files",
		},
	}
}

// subcommandAction handles the action for a given subcommand.
func subcommandAction(_ context.Context, command *cli.Command,
	generator cmd.Generator, storer cmd.Storer) error {
	privKeyPath := command.String("private")
	pubKeyPath := command.String("public")
	force := command.Bool("force")

	var errs error

	if privKeyPath == "" {
		var err error
		prompt := cmd.Prompt[string]{
			Prompt: "Private key: ",
			Validate: func(input string) (res string, err error) {
				input = strings.TrimSpace(input)
				if input == "" {
					err = fmt.Errorf("private key required")
					return
				}
				res = input
				return
			},
		}
		privKeyPath, err = prompt.Run()
		if err != nil {
			return cmd.FinalizeAction(err)
		}
	}
	if !force {
		err := protectFile(privKeyPath)
		if err != nil {
			return cmd.FinalizeAction(err)
		}
	}

	if pubKeyPath == "" {
		var err error
		prompt := cmd.Prompt[string]{
			Prompt: "Public key: ",
			Validate: func(input string) (res string, err error) {
				input = strings.TrimSpace(input)
				if input == "" {
					err = fmt.Errorf("public key required")
					return
				}
				if input == privKeyPath {
					err = fmt.Errorf("public key cannot be the same file as private key")
					return
				}
				res = input
				return
			},
		}
		pubKeyPath, err = prompt.Run()
		if err != nil {
			return cmd.FinalizeAction(err)
		}
	}
	if !force {
		err := protectFile(pubKeyPath)
		if err != nil {
			return cmd.FinalizeAction(err)
		}
	}

	var password string
	if p, set := os.LookupEnv("PASS"); set {
		password = p
	} else {
		prompt := cmd.Prompt[string]{
			Prompt:   "Password: ",
			Password: true,
			Validate: func(input string) (res string, err error) {
				if strings.HasSuffix(input, "\r\n") {
					res = strings.TrimSuffix(input, "\r\n")
				} else {
					res = strings.TrimSuffix(input, "\n")
				}
				return
			},
		}
		p, err := prompt.Run()
		if err != nil {
			return cmd.FinalizeAction(err)
		}
		password = p
	}

	color.Print("Generating keypair ")
	pubKey, privKey, err := generator(password)
	if err != nil {
		errs = cmd.JoinErrors(errs, err)
		color.Println(cmd.MsgStatusError)
		return cmd.FinalizeAction(errs)
	} else {
		color.Println(cmd.MsgStatusOK)
	}

	color.Print("Storing private key ")
	err = storer(true, privKeyPath, privKey)
	if err != nil {
		errs = cmd.JoinErrors(errs, err)
		color.Println(cmd.MsgStatusError)
	} else {
		color.Println(cmd.MsgStatusOK)
	}

	color.Print("Storing public key ")
	err = storer(false, pubKeyPath, pubKey)
	if err != nil {
		errs = cmd.JoinErrors(errs, err)
		color.Println(cmd.MsgStatusError)
	} else {
		color.Println(cmd.MsgStatusOK)
	}

	return cmd.FinalizeAction(errs)
}

// protectFile checks for existing file at path.
// If a file exist at path, the user will be prompted to confirm the overwrite.
// protectFile returns an error on rejection, and nil on confirmation or empty path.
func protectFile(path string) error {
	info, err := os.Stat(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("cannot check %s: %w", path, err)
	} else if !errors.Is(err, os.ErrNotExist) {
		prompt := cmd.Prompt[bool]{
			Prompt: color.Yellow + "Overwrite " + path + "? (y/n) " + color.Reset,
			Validate: func(input string) (bool, error) {
				input = strings.ToLower(strings.TrimSpace(input))
				if input == "y" {
					return true, nil
				} else {
					return false, nil
				}
			},
		}
		overwrite, err := prompt.Run()
		if err != nil {
			return cmd.NoReprompt(err)
		}
		if overwrite {
			if info.IsDir() {
				return fmt.Errorf("cannot overwrite directory %s", path)
			} else {
				return nil
			}
		} else {
			return cmd.NoReprompt(fmt.Errorf("existing file exist at %s", path))
		}
	}
	return nil
}
