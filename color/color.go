package color

import (
	"fmt"
	"os"
	"regexp"

	"github.com/mattn/go-isatty"
)

const (
	Black = "\033[30m"
	White = "\033[97m"

	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[90m"

	LightRed     = "\033[91m"
	LightGreen   = "\033[92m"
	LightYellow  = "\033[93m"
	LightBlue    = "\033[94m"
	LightMagenta = "\033[95m"
	LightCyan    = "\033[96m"
	LightGray    = "\033[37m"

	Reset        = "\033[0m"
	Bold         = "\033[1m"
	Underline    = "\033[4m"
	NoUnderline  = "\033[24m"
	ReverseText  = "\033[7m"
	PositiveText = "\033[27m"

	EscStr  = "\033"
	EscChar = '\033'
)

var pattern = regexp.MustCompile(`(?:\033\[[0-9]{1,2}m)`)

func Print(a ...any) (n int, err error) {
	s := nonTtyClean(fmt.Sprint(a...))

	return fmt.Print(s)
}

func Printf(format string, a ...any) (n int, err error) {
	s := nonTtyClean(fmt.Sprintf(format, a...))

	return fmt.Print(s)
}

func Println(a ...any) (n int, err error) {
	s := nonTtyClean(fmt.Sprintln(a...))

	return fmt.Print(s)
}

func nonTtyClean(s string) string {
	if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return s
	}

	return CleanString(s)
}

// Clean strips all formatting from b.
func Clean(b []byte) []byte {
	return pattern.ReplaceAll(b, []byte{})
}

// CleanString strips all formatting from s.
func CleanString(s string) string {
	return pattern.ReplaceAllString(s, "")
}
