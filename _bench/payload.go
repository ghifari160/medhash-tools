// MedHash Tools
// Copyright (c) 2023 GHIFARI160
// MIT License

package main

import (
	"crypto/rand"
	"fmt"
	"os"
)

const maxBuffer = 1 * 1024 * 1024 * 1024

// GenPayload generates payload of the specified size.
func GenPayload(path string, size Size) error {
	var buf []byte
	var counter int64

	if size > maxBuffer {
		buf = make([]byte, maxBuffer)
	} else {
		buf = make([]byte, size)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create payload: %w", err)
	}

	n, err := rand.Read(buf)
	if err != nil {
		return fmt.Errorf("cannot generate payload: %w", err)
	}

	if counter+int64(n) > int64(size) {
		_, err = f.Write(buf[:int64(size)-counter])
	} else {
		_, err = f.Write(buf)
	}
	if err != nil {
		return fmt.Errorf("cannot write to payload: %w", err)
	}

	return f.Close()
}

// TempDir creates a new temporary directory and returns the absolute path.
// If one cannot be created, TempDir panics.
func TempDir() string {
	path, err := os.MkdirTemp("", "medhash-tools_benchmark")
	if err != nil {
		panic(fmt.Errorf("cannot create temporary directory: %w", err))
	}

	return path
}
