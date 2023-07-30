package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	var args Args

	args.MustParse(os.Args)

	if args.Count < 1 || args.PayloadSize < 1 || len(args.Cmd) < 1 {
		args.PrintHelp()
		os.Exit(1)
	}

	dir := TempDir()
	defer func() {
		fmt.Println("Cleaning up")

		err := os.RemoveAll(dir)
		if err != nil {
			fmt.Printf("Cannot clean up: %v\n", err)
		}
	}()

	if args.Report != "" {
		fmt.Println("Executing in report mode!")

		err := os.MkdirAll(args.Report, 0777)
		if err != nil {
			fmt.Printf("Cannot create report directory: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("Generating payload to %s\n", dir)

	err := GenPayload(filepath.Join(dir, "payload"), args.PayloadSize)
	if err != nil {
		fmt.Printf("Cannot generate payload: %v\n", err)
		os.Exit(1)
	}

	cases := make(map[string]Result)

	for _, cmd := range args.Cmd {
		var report *os.File
		var printer io.Writer

		cmd.Args = append(cmd.Args, dir)

		if args.Report != "" {
			report, err := os.Create(filepath.Join(args.Report, filepath.Base(cmd.Cmd)))
			if err != nil {
				fmt.Printf("Cannot create report file: %v\n", err)
				os.Exit(1)
			}
			printer = io.MultiWriter(report, os.Stdout)
		}

		fmt.Fprintf(printer, "Benchmarking %s with arguments [ %s ]\n", cmd.Cmd,
			strings.Join(cmd.Args, ", "))

		results := make([]Result, 0)

		for i := 0; i < args.Count; i++ {
			fmt.Printf("%d of %d\n", i+1, args.Count)
			results = append(results, bench(cmd.Cmd, cmd.Args, args.PayloadSize))
		}

		var result Result

		for _, r := range results {
			result.Duration += r.Duration
			result.Rate += r.Rate
		}

		result.Iteration = len(results)
		result.Duration /= time.Duration(len(results))
		result.Rate /= float64(len(results))
		result.PayloadSize = args.PayloadSize

		cases[filepath.Base(cmd.Cmd)] = result
		fmt.Fprint(printer, result.Report())

		if report != nil {
			err = report.Close()
			if err != nil {
				fmt.Printf("Cannot write to report file: %v\n", err)
				os.Exit(1)
			}
		}
	}

	if args.Report != "" {
		tbl := make([][]string, 1)

		tbl[0] = []string{"Program", "Duration (s)", "Rate (GiB/s)", "Payload Size"}

		for label, result := range cases {
			tbl = append(tbl, []string{
				label,
				fmt.Sprintf("%.4f", result.Duration.Seconds()),
				fmt.Sprintf("%.4f", result.Rate),
				result.PayloadSize.String(),
			})
		}

		summary := Summary(tbl)

		err := os.WriteFile(filepath.Join(args.Report, "summary.md"), []byte(summary), 0644)
		if err != nil {
			fmt.Printf("Cannot write summary: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Report stored to %s\n", args.Report)
	}
}

// bench executes the benchmark.
func bench(cmdName string, cmdArgs []string, payloadSize Size) (res Result) {
	cmd := exec.Command(cmdName, cmdArgs...)

	start := time.Now()
	out, err := cmd.CombinedOutput()
	stop := time.Now()

	if err != nil {
		fmt.Println(string(out))

		panic(err)
	}

	res.Duration = stop.Sub(start)
	res.Rate = payloadSize.GiB() / res.Duration.Seconds()
	res.PayloadSize = payloadSize

	return
}
