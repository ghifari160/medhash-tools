package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/ghifari160/medhash-tools/cmd"
)

const Name = "MedHash Tools"
const Version = "0.6.0"

func main() {
	var args struct {
		Gen *cmd.Gen        `arg:"subcommand:gen" help:"generate MedHash Manifest"`
		Ver *cmd.GenericCmd `arg:"subcommand:version" help:"print tool version"`
	}

	p := arg.MustParse(&args)

	if p.Subcommand() == nil {
		p.WriteUsage(os.Stdout)
		os.Exit(1)
	}

	printHeader()

	var c cmd.Command

	switch {
	case args.Gen != nil:
		c = args.Gen

	case args.Ver != nil:
		c = new(cmd.GenericCmd)
	}

	os.Exit(c.Execute())
}

func printHeader() {
	fmt.Printf("%s v%s\n", Name, Version)
}
