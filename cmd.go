package main

// Cmd is a generic wrapper for all subcommands.
type Cmd interface {
	Execute() int
}

// GenericCmd is a generic subcommand.
type GenericCmd struct{}

func (c GenericCmd) Execute() (status int) { return }
