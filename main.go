package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
)

const Name = "MedHash Tools"
const Version = "0.6.0"

func main() {
	var args struct {
		Ver *GenericCmd `arg:"subcommand:version"`
	}

	p := arg.MustParse(&args)

	if p.Subcommand() == nil {
		p.WriteUsage(os.Stdout)
		os.Exit(1)
	}

	printHeader()

	var cmd Cmd

	switch {
	case args.Ver != nil:
		cmd = new(GenericCmd)
	}

	os.Exit(cmd.Execute())
}

func printHeader() {
	fmt.Printf("%s v%s\n", Name, Version)
}
