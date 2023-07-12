// MedHash Tools
// Copyright (c) 2023 GHIFARI160
// MIT License

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/alexflint/go-arg"
)

func main() {
	var args struct {
		Count       int      `arg:"--count,-q" default:"1" help:"Iteration count"`
		Cmd         string   `arg:"--command,-c,required" help:"Command to becnhmark"`
		Args        []string `arg:"--args,-a" help:"Arguments for the command"`
		PayloadSize int      `arg:"--size,-s" default:"4" help:"Payload size in GiB"`
	}

	arg.MustParse(&args)

	dir := TempDir()
	defer func() {
		fmt.Println("Cleaning up")
		os.RemoveAll(dir)
	}()

	payloadSize := Size(args.PayloadSize) * GiB

	fmt.Println("Generating payload")

	err := GenPayload(filepath.Join(dir, "payload"), payloadSize)
	if err != nil {
		panic(err)
	}

	var cmdArgs []string

	if len(args.Args) > 0 {
		cmdArgs = args.Args
	} else {
		cmdArgs = make([]string, 0)
	}

	cmdArgs = append(cmdArgs, dir)

	fmt.Printf("Benchmarking %s\n", args.Cmd)

	var result Result
	results := make([]Result, 0)

	for i := 0; i < args.Count; i++ {
		fmt.Printf("%d of %d\n", i+1, args.Count)
		results = append(results, bench(args.Cmd, cmdArgs, payloadSize))
	}

	for _, r := range results {
		result.Duration += r.Duration
		result.Rate += r.Rate
	}

	result.Duration /= time.Duration(len(results))
	result.Rate /= float64(len(results))

	fmt.Printf("Duration: %.4f s\n", result.Duration.Seconds())
	fmt.Printf("Rate:     %.4f GiB/s\n", result.Rate)
	fmt.Printf("Size:     %s\n", payloadSize)
}

// Result holds the result of the benchmark.
type Result struct {
	Duration time.Duration
	Rate     float64
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

	return
}
