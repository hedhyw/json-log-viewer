package fileinput

import (
	"context"
	"io"
	"os"
)

// FileInput is the source that represents the file.
type FileInput struct {
	fileName string
}

// New initializes a new FileInput with the given file.
func New(fileName string) FileInput {
	return FileInput{
		fileName: fileName,
	}
}

// ReadCloser opens the file. Call Close after usage.
func (s FileInput) ReadCloser(context.Context) (io.ReadCloser, error) {
	return os.Open(s.fileName)
}

// String implements fmt.Stringer.
func (s FileInput) String() string {
	return s.fileName
}
