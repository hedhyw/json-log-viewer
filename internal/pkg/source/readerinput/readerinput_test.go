package readerinput_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"testing/iotest"
	"time"

	"github.com/hedhyw/json-log-viewer/assets"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source/readerinput"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReaderInput(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expected := assets.ExampleJSONLog()

	t.Run("ReadCloser", func(t *testing.T) {
		t.Parallel()

		input := readerinput.New(bytes.NewReader(expected), time.Minute)

		readCloser, err := input.ReadCloser(ctx)
		require.NoError(t, err)

		t.Cleanup(func() { assert.NoError(t, readCloser.Close()) })

		actual, err := io.ReadAll(readCloser)
		require.NoError(t, err)

		assert.Equal(t, bytes.TrimSpace(expected), bytes.TrimSpace(actual))
	})

	t.Run("ReadCloser_twice", func(t *testing.T) {
		t.Parallel()

		input := readerinput.New(bytes.NewReader(expected), time.Minute)

		for range 2 {
			readCloser, err := input.ReadCloser(ctx)
			require.NoError(t, err)

			t.Cleanup(func() { assert.NoError(t, readCloser.Close()) })

			actual, err := io.ReadAll(readCloser)
			require.NoError(t, err)

			assert.Equal(t, bytes.TrimSpace(expected), bytes.TrimSpace(actual))
		}
	})

	t.Run("ReadCloser_error", func(t *testing.T) {
		t.Parallel()

		// nolint: err113 // Test error.
		errReader := errors.New(t.Name())

		input := readerinput.New(iotest.ErrReader(errReader), time.Minute)

		_, err := input.ReadCloser(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, errReader)

		_, err = input.ReadCloser(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, errReader)
	})

	t.Run("ReadCloser_wait", func(t *testing.T) {
		t.Parallel()

		const (
			lineFirst  = "line first\n"
			lineSecond = "line second\n"

			timeout = 200 * time.Millisecond
		)

		pipeReader, pipeWriter := io.Pipe()

		t.Cleanup(func() { assert.NoError(t, pipeReader.Close()) })
		t.Cleanup(func() { assert.NoError(t, pipeWriter.Close()) })

		input := readerinput.New(pipeReader, timeout)

		_, err := pipeWriter.Write([]byte(lineFirst))
		require.NoError(t, err)

		readCloser, err := input.ReadCloser(ctx)
		require.NoError(t, err)

		actual, err := io.ReadAll(readCloser)
		require.NoError(t, err)

		assert.Equal(t, lineFirst, string(actual))

		_, err = pipeWriter.Write([]byte(lineSecond))
		require.NoError(t, err)

		readCloser, err = input.ReadCloser(ctx)
		require.NoError(t, err)

		actual, err = io.ReadAll(readCloser)
		require.NoError(t, err)

		assert.Equal(t, lineFirst+lineSecond, string(actual))
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		input := readerinput.New(bytes.NewReader(nil), time.Minute)

		assert.Equal(t, "-", input.String())
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		input := readerinput.New(bytes.NewReader(nil), time.Minute)

		readCloser, err := input.ReadCloser(ctx)
		require.NoError(t, err)

		t.Cleanup(func() { assert.NoError(t, readCloser.Close()) })

		actual, err := io.ReadAll(readCloser)
		require.NoError(t, err)

		assert.Empty(t, actual)
	})
}
