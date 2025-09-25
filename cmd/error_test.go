package cmd_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ghifari160/medhash-tools/cmd"
	"github.com/stretchr/testify/assert"
)

func TestJoinErrors(t *testing.T) {
	cases := []struct {
		id   string
		errs []error
	}{
		{
			id: "single",
			errs: []error{
				errors.New("1"),
				errors.New("2"),
				errors.New("3"),
			},
		},
		{
			id: "nested/single",
			errs: []error{
				errors.New("1"),
				errors.New("2"),
				fmt.Errorf("3%w", errors.New("a")),
			},
		},
		{
			id: "nested/multiple",
			errs: []error{
				errors.New("1"),
				errors.New("2"),
				cmd.JoinErrors(
					errors.New("3a"),
					errors.New("3b"),
					errors.New("3c"),
				),
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.id, func(t *testing.T) {
			testJoinErrors(t, testCase.errs)
		})
	}
}

func TestUnwrapJoinedErrors(t *testing.T) {
	cases := []struct {
		id   string
		errs []error
	}{
		{
			id: "single",
			errs: []error{
				errors.New("1"),
				errors.New("2"),
				errors.New("3"),
			},
		},
		{
			id: "nested",
			errs: []error{
				errors.New("1"),
				errors.New("2"),
				cmd.JoinErrors(
					errors.New("3a"),
					errors.New("3b"),
					errors.New("3c"),
				),
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.id, func(t *testing.T) {
			testUnwrapJoinedErrors(t, testCase.errs)
		})
	}
}

func testJoinErrors(t *testing.T, errs []error) {
	assert := assert.New(t)

	expected := countErrors(t, errs...)

	err := cmd.JoinErrors(errs...)
	actual := countErrors(t, err)

	assert.Equal(expected, actual)
}

func testUnwrapJoinedErrors(t *testing.T, errs []error) {
	assert := assert.New(t)

	expected := countErrors(t, errs...)

	unwrapped := cmd.UnwrapJoinedErrors(cmd.JoinErrors(errs...))
	actual := countErrors(t, unwrapped...)

	assert.Equal(expected, actual)
}

func countErrors(t *testing.T, errs ...error) int {
	t.Helper()
	var expected int
	for _, err := range errs {
		if joinedErrs, ok := err.(interface{ Unwrap() []error }); ok {
			expected += countErrors(t, joinedErrs.Unwrap()...)
		} else {
			expected++
		}
	}
	return expected
}
