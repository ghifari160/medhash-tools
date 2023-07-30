package main

import (
	"fmt"
	"time"
)

// Result holds the result of the benchmark.
type Result struct {
	Iteration   int
	Duration    time.Duration
	Rate        float64
	PayloadSize Size
}

// Report prints a report from this Result.
func (r Result) Report() string {
	var report string

	if r.Iteration > 0 {
		report += fmt.Sprintf("Iteration: %d\n", r.Iteration)
	}

	report += fmt.Sprintf("Duration:  %.4f\n", r.Duration.Seconds())
	report += fmt.Sprintf("Rate:      %.4f\n", r.Rate)
	report += fmt.Sprintf("Size:      %s\n", r.PayloadSize)

	return report
}
