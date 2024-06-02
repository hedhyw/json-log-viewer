package source

import (
	"context"
	"fmt"
	"io"
)

// Input returns the getter of read-closer for the given input source.
type Input interface {
	// ReadCloser returns a reader from the input. Call `Close` after usage.
	ReadCloser(ctx context.Context) (io.ReadCloser, error)

	fmt.Stringer
}
