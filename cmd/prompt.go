package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/ghifari160/medhash-tools/color"
)

// Prompt is a simple prompter with validation.
type Prompt[R any] struct {
	Prompt string
	// Validate validates input and transforms it to R.
	// If Validate returns an error, the user will be re-prompted.
	// Wrap the error with NoReprompt to prevent re-prompts.
	Validate func(input string) (R, error)
}

// Run begins to prompt the user and validates the received input.
func (p *Prompt[T]) Run() (result T, err error) {
	if p.Validate == nil {
		err = fmt.Errorf("p.Validate must not be nil")
		return
	}
	reader := bufio.NewReader(os.Stdin)
	return p.subRun(reader)
}

// subRun prompts the user and validates the input.
// If p.Validate returns a non-NoReprompt error, subRun calls itself recursively to reprompt the
// user.
func (p *Prompt[T]) subRun(reader *bufio.Reader) (res T, err error) {
	color.Print(p.Prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	res, err = p.Validate(input)
	if err != nil && !errors.Is(err, noRepromptErr{}) {
		color.Println(color.LightRed + UpperCaseFirst(err.Error()) + color.Reset)
		return p.subRun(reader)
	}
	return
}

// NoReprompt wraps err to prevent Prompt from re-prompting the user.
func NoReprompt(err error) error {
	if err != nil {
		return noRepromptErr{err: err}
	} else {
		return nil
	}
}

// noRepromptErr wraps an error in a way that prevents Prompt from re-prompting the user.
type noRepromptErr struct {
	err error
}

func (e noRepromptErr) Error() string {
	return e.err.Error()
}

func (e noRepromptErr) Is(tgt error) bool {
	_, ok := tgt.(noRepromptErr)
	return ok
}

func (e noRepromptErr) Unwrap() error {
	return e.err
}
