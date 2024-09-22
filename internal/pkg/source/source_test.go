package source_test

import (
	"bytes"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hedhyw/json-log-viewer/assets"
	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
	"github.com/hedhyw/json-log-viewer/internal/pkg/tests"
	"github.com/hedhyw/semerr/pkg/v1/semerr"
)

func TestParseLogEntries(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		reader := bytes.NewReader(assets.ExampleJSONLog())

		source, err := source.Reader(reader, config.GetDefaultConfig())
		require.NoError(t, err)

		t.Cleanup(func() { assert.NoError(t, source.Close()) })

		logEntries, err := source.ParseLogEntries()

		require.NoError(t, err)
		assert.NotEmpty(t, logEntries)
	})

	t.Run("large_line", func(t *testing.T) {
		t.Parallel()

		longLine := strings.Repeat("1", 2*1024*1024)

		reader := strings.NewReader(longLine)

		source, err := source.Reader(reader, config.GetDefaultConfig())
		require.NoError(t, err)

		t.Cleanup(func() { assert.NoError(t, source.Close()) })

		logEntries, err := source.ParseLogEntries()

		require.NoError(t, err)
		assert.NotEmpty(t, logEntries)
	})

	t.Run("failed", func(t *testing.T) {
		t.Parallel()

		reader := iotest.ErrReader(semerr.Error("test"))

		source, err := source.Reader(reader, config.GetDefaultConfig())
		require.NoError(t, err)

		t.Cleanup(func() { assert.NoError(t, source.Close()) })

		_, err = source.ParseLogEntries()
		require.Error(t, err)
	})
}

func TestParseLogEntriesFromReaderLimited(t *testing.T) {
	t.Parallel()

	content := `{}`

	cfg := config.GetDefaultConfig()
	cfg.MaxFileSizeBytes = 1

	reader := strings.NewReader(content)
	is, err := source.Reader(reader, cfg)
	require.NoError(t, err)
	defer is.Close()
	logEntries, err := is.ParseLogEntries()
	require.NoError(t, err)

	require.Empty(t, logEntries.Entries)
}

func TestRow(t *testing.T) {
	t.Parallel()

	entry := t.Name() + "\n"

	input := bytes.NewReader([]byte(entry))

	cfg := config.GetDefaultConfig()

	source, err := source.Reader(input, cfg)
	require.NoError(t, err)

	t.Cleanup(func() { assert.NoError(t, source.Close()) })

	lazyEntries, err := source.ParseLogEntries()
	require.NoError(t, err)

	assert.Equal(t, 1, lazyEntries.Len())

	row := lazyEntries.Row(cfg, 0)
	assert.Contains(t, row, entry)
}

func TestFile(t *testing.T) {
	t.Parallel()

	cfg := config.GetDefaultConfig()
	fileName := tests.RequireCreateFile(t, []byte(t.Name()+"\n"))

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		source, err := source.File(fileName, cfg)
		require.NoError(t, err)

		assert.True(t, source.CanFollow())
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()

		_, err := source.File(fileName+"-not-found", cfg)
		require.Error(t, err)
	})
}
