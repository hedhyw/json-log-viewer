//go:build !mock_stdin

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
)

func getStdinReader(defaultInput fs.File) (io.Reader, error) {
	stat, err := defaultInput.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}

	if stat.Mode()&os.ModeCharDevice != 0 {
		return bytes.NewReader(nil), nil
	}
	return defaultInput, nil
}
