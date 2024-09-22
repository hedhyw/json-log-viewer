package source_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
	"github.com/hedhyw/json-log-viewer/internal/pkg/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDelay = source.RefreshInterval * 2

func TestStartStreamingEndOfFile(t *testing.T) {
	t.Parallel()

	entry := t.Name() + "\n"

	input := bytes.NewReader([]byte(entry))

	cfg := config.GetDefaultConfig()

	inputSource, err := source.Reader(input, cfg)
	require.NoError(t, err)

	t.Cleanup(func() { assert.NoError(t, inputSource.Close()) })

	ctx, cancel := context.WithCancel(tests.Context(t))

	entries := make(chan source.LazyLogEntries)

	inputSource.StartStreaming(ctx, func(msg source.LazyLogEntries, err error) {
		require.NoError(t, err)

		select {
		case entries <- msg:
		case <-ctx.Done():
		}
	})

	select {
	case msg := <-entries:
		cancel()

		require.Equal(t, msg.Len(), 1)
		assert.Contains(t, msg.Row(cfg, 0), entry)
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}

func TestStartStreamingUpdates(t *testing.T) {
	t.Parallel()

	entry := t.Name() + "\n"

	pipeReader, pipeWriter := io.Pipe()

	t.Cleanup(func() {
		assert.NoError(t, pipeReader.Close())
		assert.NoError(t, pipeWriter.Close())
	})

	cfg := config.GetDefaultConfig()

	inputSource, err := source.Reader(pipeReader, cfg)
	require.NoError(t, err)

	t.Cleanup(func() { assert.NoError(t, inputSource.Close()) })

	ctx := tests.Context(t)

	entries := make(chan source.LazyLogEntries)

	inputSource.StartStreaming(ctx, func(msg source.LazyLogEntries, err error) {
		if msg.Len() == 0 {
			return
		}

		select {
		case entries <- msg:
		case <-ctx.Done():
		}
	})

	for i := range 2 {
		_, err = fmt.Fprintln(pipeWriter, entry)
		require.NoError(t, err)

		select {
		case msg := <-entries:
			require.Equalf(t, i+1, msg.Len(), "iteration %d", i)
			assert.Containsf(t, msg.Row(cfg, 0), entry, "iteration %d", i)
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		}

		select {
		case <-time.After(testDelay):
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		}
	}
}

func TestStartStreamingContextClosed(t *testing.T) {
	t.Parallel()

	pipeReader, pipeWriter := io.Pipe()

	t.Cleanup(func() {
		assert.NoError(t, pipeReader.Close())
		assert.NoError(t, pipeWriter.Close())
	})

	cfg := config.GetDefaultConfig()

	inputSource, err := source.Reader(pipeReader, cfg)
	require.NoError(t, err)

	t.Cleanup(func() { assert.NoError(t, inputSource.Close()) })

	ctx, cancel := context.WithCancel(tests.Context(t))
	defer cancel()

	inputSource.StartStreaming(ctx, func(source.LazyLogEntries, error) {})

	cancel()

	<-time.After(2 * source.RefreshInterval)
}

func TestStartStreamingFromFile(t *testing.T) {
	t.Parallel()

	entry := t.Name() + "\n"

	fileName := tests.RequireCreateFile(t, []byte(""))

	cfg := config.GetDefaultConfig()

	inputSource, err := source.File(fileName, cfg)
	require.NoError(t, err)

	t.Cleanup(func() { assert.NoError(t, inputSource.Close()) })

	ctx := tests.Context(t)

	entries := make(chan source.LazyLogEntries)

	inputSource.StartStreaming(ctx, func(msg source.LazyLogEntries, err error) {
		if msg.Len() == 0 {
			return
		}

		select {
		case entries <- msg:
		case <-ctx.Done():
		}
	})

	file, err := os.OpenFile(fileName, os.O_WRONLY, os.ModePerm)
	require.NoError(t, err)

	t.Cleanup(func() { assert.NoError(t, file.Close()) })

	for i := range 2 {
		_, err = fmt.Fprintln(file, entry)
		require.NoError(t, err)

		require.NoError(t, file.Sync())

		select {
		case msg := <-entries:
			require.Equalf(t, i+1, msg.Len(), "iteration %d", i)
			assert.Containsf(t, msg.Row(cfg, 0), entry, "iteration %d", i)
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		}

		select {
		case <-time.After(testDelay):
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		}
	}
}
