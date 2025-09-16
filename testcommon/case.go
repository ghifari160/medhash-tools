package testcommon

import "testing"

// TestFn tests a test case.
type TestFn func(t *testing.T, alg string, opts ...Options)

// TestCase represents a given test case.
type TestCase struct {
	ID      string
	Alg     string
	Options []Options
}

func (testCase TestCase) Run(t *testing.T, testFn TestFn) {
	t.Run(testCase.ID, func(t *testing.T) {
		testFn(t, testCase.Alg, testCase.Options...)
	})
}

// Case constructs a new TestCase.
func Case(id, alg string, opts ...Options) TestCase {
	return TestCase{id, alg, opts}
}

// RunCases runs all TestCase in cases.
func RunCases(t *testing.T, testFn TestFn, cases []TestCase) {
	for _, testCase := range cases {
		testCase.Run(t, testFn)
	}
}
