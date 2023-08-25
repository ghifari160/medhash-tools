package cmd

import "github.com/ghifari160/medhash-tools/color"

const (
	MsgStatusError   = color.Red + "ERROR" + color.Reset
	MsgStatusOK      = color.Green + "OK" + color.Reset
	MsgStatusSkipped = color.Yellow + "SKIPPED" + color.Reset
	MsgFinalError    = color.Red + "Error!" + color.Reset
	MsgFinalDone     = color.Green + "Done!" + color.Reset
)

// Command is a generic wrapper for all subcommands.
type Command interface {
	Execute() (status int)
}

// GenericCmd is a generic subcommand.
type GenericCmd struct{}

func (c GenericCmd) Execute() (status int) { return }

// CmdConfig is a common set of parameters for the commands.
type CmdConfig struct {
	Default bool `arg:"--default,-d" default:"true" help:"use default preset"`
	All     bool `arg:"--all,-a" help:"use all algorithms"`

	XXH3   bool `arg:"--xxh3" help:"use XXH3"`
	SHA3   bool `arg:"--sha3" help:"use SHA3-256"`
	SHA256 bool `arg:"--sha256" help:"use SHA256"`
	SHA1   bool `arg:"--sha1" help:"use SHA1"`
	MD5    bool `arg:"--md5" help:"use MD5"`
}
