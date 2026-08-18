package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
)

// stdinInput is a log source that is read from the standard input.
type stdinInput struct {
	// Reader of the standard input.
	Reader io.Reader
	// IsPipe is true if another command writes to the standard input.
	IsPipe bool
}

func getStdinReader(defaultInput fs.File) (stdinInput, error) {
	stat, err := defaultInput.Stat()
	if err != nil {
		return stdinInput{}, fmt.Errorf("stat: %w", err)
	}

	if stat.Mode()&os.ModeCharDevice != 0 {
		return stdinInput{Reader: bytes.NewReader(nil)}, nil
	}

	return stdinInput{
		Reader: defaultInput,
		IsPipe: stat.Mode()&os.ModeNamedPipe != 0,
	}, nil
}
