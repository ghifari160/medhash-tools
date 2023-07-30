package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

var ErrHelp = errors.New("print help message")

// Args stores the parameters for this program.
type Args struct {
	Count       int
	PayloadSize Size
	Report      string

	Cmd []Cmd
}

// Cmd stores information about the command to benchmark.
type Cmd struct {
	Label string
	Cmd   string
	Args  []string
}

// PrintHelp prints the help message.
func (ar *Args) PrintHelp() {
	fmt.Printf("Usage: %s --count COUNT --size SIZE --cmd CMD [ARGS...] [--cmd CMD [ARGS...]...]\n", "_bench")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --count, -q COUNT        iteration count")
	fmt.Println("  --size, -s SIZE          payload size. Use prefixes \"K\" for KiB, \"M\" for MiB, and \"G\" for GiB")
	fmt.Println("  --report, -r REPORT      path to store report")
	fmt.Println("  --cmd, -c CMD [ARGS...]  command to benchmark")
	fmt.Println("  --help, -h               display this help and exit")
}

// MustParse attempts to parse arguments.
// Otherwise, MustParse calls PrintHelp.
func (ar *Args) MustParse(args []string) {
	err := ar.Parse(args)
	if err != nil {
		ar.PrintHelp()
		os.Exit(1)
	}
}

// Parse attempts to parse argument.
// Otherwise, Parse returns an error.
func (ar *Args) Parse(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("missing arguments")
	}

	err := ar.parse(args[1:])
	if errors.Is(err, ErrHelp) {
		ar.PrintHelp()
		os.Exit(1)
	}

	return err
}

// parse attempts to parse arguments.
func (ar *Args) parse(args []string) error {
	var c Cmd

	ar.Cmd = make([]Cmd, 0)

	for i := 0; i < len(args); i++ {
		var op string

		a := args[i]

		if len(a) < 1 {
			continue
		}

		if len(a) > 2 && a[:2] == "--" {
			switch a[2:] {
			case "help":
				op = "help"

			case "count":
				op = "count"

			case "size":
				op = "size"

			case "report":
				op = "report"

			case "cmd":
				op = "cmd"

			case "label":
				op = "label"
			}
		} else if len(a) > 1 && a[:1] == "-" {
			switch a[1:] {
			case "h":
				op = "help"

			case "q":
				op = "count"

			case "s":
				op = "size"

			case "r":
				op = "report"

			case "c":
				op = "cmd"

			case "l":
				op = "label"
			}
		}

		switch op {
		case "help":
			return ErrHelp

		case "count":
			q, err := strconv.Atoi(args[i+1])
			if err != nil {
				return err
			}
			ar.Count = q
			i++

		case "size":
			s, err := ar.parseSize(args[i+1])
			if err != nil {
				return err
			}
			ar.PayloadSize = s
			i++

		case "report":
			ar.Report = args[i+1]
			i++

		case "label":
			if c.Cmd != "" && c.Label != "" {
				ar.Cmd = append(ar.Cmd, c)
				c = Cmd{
					Args: make([]string, 0),
				}
			}
			c.Label = args[i+1]
			i++

		case "cmd":
			if c.Cmd != "" && c.Label != "" {
				ar.Cmd = append(ar.Cmd, c)
				c = Cmd{
					Args: make([]string, 0),
				}
			}
			c.Cmd = args[i+1]
			i++

		default:
			c.Args = append(c.Args, a)
		}
	}

	if c.Cmd != "" {
		ar.Cmd = append(ar.Cmd, c)
	}

	return nil
}

// parseSize parses a size string into a Size.
// It attempts to interpret the unit suffix as one of either "K", "M", or "G".
// If the unit suffix cannot be determined, parseSize assumes the size is in bytes.
func (ar *Args) parseSize(size string) (Size, error) {
	var s Size
	var operandStr string

	if len(size) < 1 {
		return 0, fmt.Errorf("unknown size: %s", size)
	}

	unit := size[len(size)-1:]

	switch unit {
	case "K":
		s = KiB

	case "M":
		s = MiB

	case "G":
		s = GiB

	default:
		s = 1
		operandStr = size
	}

	if operandStr == "" {
		operandStr = size[:len(size)-1]
	}

	operand, err := strconv.Atoi(operandStr)
	if err != nil {
		return 0, err
	}
	s *= Size(operand)

	return s, nil
}
