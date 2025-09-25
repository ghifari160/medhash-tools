package cmd

import (
	"context"

	"github.com/ghifari160/medhash-tools/color"
	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/urfave/cli/v3"
)

func init() {
	RegisterCmd(&cli.Command{
		Name:  "version",
		Usage: "print tool version",
		Action: func(ctx context.Context, c *cli.Command) error {
			return nil
		},
	})
}

var commands []*cli.Command

// RegisterCmd registers a command.
func RegisterCmd(cmd *cli.Command) {
	commands = append(commands, cmd)
}

// Commands returns all registered commands.
// Note that each command package may need to be anonymously imported.
func Commands() []*cli.Command {
	return commands
}

// HashAlgs returns the appropriate flags for all supported hashing algorithms.
func HashAlgs() []cli.Flag {
	return []cli.Flag{
		simpleBoolFlag("xxh3", "use XXH3"),
		simpleBoolFlag("sha512", "use SHA512"),
		simpleBoolFlag("sha3", "use SHA3"),
		simpleBoolFlag("sha256", "use SHA256"),
		simpleBoolFlag("sha1", "use SHA1"),
		simpleBoolFlag("md5", "use MD5"),
	}
}

// ConfigFromFlags returns config based on the flags of command.
// The flags checked are the ones HashAlgs returns.
func ConfigFromFlags(command *cli.Command) medhash.Config {
	var config medhash.Config

	if command.Bool("all") {
		return medhash.AllConfig
	}

	defaultConf := true
	if command.Bool("xxh3") {
		defaultConf = false
		config.XXH3 = true
	}
	if command.Bool("sha512") {
		defaultConf = false
		config.SHA512 = true
	}
	if command.Bool("sha3") {
		defaultConf = false
		config.SHA3 = true
	}
	if command.Bool("sha256") {
		defaultConf = false
		config.SHA256 = true
	}
	if command.Bool("sha1") {
		defaultConf = false
		config.SHA1 = true
	}
	if command.Bool("md5") {
		defaultConf = false
		config.MD5 = true
	}

	if defaultConf {
		config = medhash.DefaultConfig
	}

	return config
}

// simpleBoolFlag returns a new cli.BoolFlag with just the Name and Usage set.
func simpleBoolFlag(name, usage string) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  name,
		Usage: usage,
	}
}

const (
	MsgStatusError   = color.Red + "ERROR" + color.Reset
	MsgStatusOK      = color.Green + "OK" + color.Reset
	MsgStatusSkipped = color.Yellow + "SKIPPED" + color.Reset
	MsgFinalError    = color.Red + "Error!" + color.Reset
	MsgFinalDone     = color.Green + "Done!" + color.Reset
)
