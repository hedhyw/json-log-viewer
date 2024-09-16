//go:build mock_stdin

package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"time"
)

func getStdinReader(defaultInput fs.File) (io.Reader, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("pipe: %w", err)
	}
	go func() {
		defer w.Close()
		for i := 0; ; i++ {
			_, err := w.Write([]byte(fmt.Sprintf(`{"message": "Line %d"}
	`, i)))
			if err != nil {
				fatalf("Write failed: %s\n", err)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()
	return r, nil
}
