package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ghifari160/medhash-tools/cmd"
	_ "github.com/ghifari160/medhash-tools/cmd/chk"
	_ "github.com/ghifari160/medhash-tools/cmd/gen"
	_ "github.com/ghifari160/medhash-tools/cmd/upgrade"
	"github.com/urfave/cli/v3"
)

const Name = "MedHash Tools"
const Version = "0.7.0"

func main() {
	root := &cli.Command{
		Name:     "medhash",
		Usage:    "Simple tool for verifying media file integrity",
		Commands: cmd.Commands(),
	}

	printHeader()
	err := root.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Printf("main: %v\n", err)
	}
}

func printHeader() {
	fmt.Printf("%s v%s\n", Name, Version)
}
