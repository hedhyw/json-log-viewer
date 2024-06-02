package readerinput

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"time"
)

// ReaderInput reads from the configured input with some timeout.
type ReaderInput struct {
	readTimeout time.Duration

	linesChan <-chan []byte
	errChan   <-chan error

	lastErr error
	content []byte
}

// New initializes ReaderInput with the given reader and timeout.
func New(
	reader io.Reader,
	timeout time.Duration,
) *ReaderInput {
	scanner := bufio.NewScanner(reader)

	linesChan := make(chan []byte)
	errChan := make(chan error, 1)

	go func() {
		for scanner.Scan() {
			linesChan <- scanner.Bytes()
		}

		if err := scanner.Err(); err != nil {
			errChan <- err
			close(errChan)
		} else {
			close(linesChan)
		}
	}()

	return &ReaderInput{
		readTimeout: timeout,

		linesChan: linesChan,
		errChan:   errChan,

		lastErr: nil,
		content: make([]byte, 0),
	}
}

// String implements fmt.Stringer.
func (s *ReaderInput) String() string {
	return "-"
}

// ReadCloser reads the content from the input.
func (s *ReaderInput) ReadCloser(ctx context.Context) (io.ReadCloser, error) {
	if s.lastErr != nil {
		return nil, s.lastErr
	}

	ctx, cancel := context.WithTimeout(ctx, s.readTimeout)
	defer cancel()

loop:
	for ctx.Err() == nil {
		select {
		case line, ok := <-s.linesChan:
			if !ok {
				break loop
			}

			s.content = append(s.content, line...)
			s.content = append(s.content, "\n"...)
		case err := <-s.errChan:
			s.lastErr = err

			return nil, s.lastErr
		case <-ctx.Done():
			break loop
		}
	}

	return io.NopCloser(bytes.NewReader(s.content)), nil
}
