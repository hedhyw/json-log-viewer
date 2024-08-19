package source_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hedhyw/json-log-viewer/assets"
	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

func TestLoadLogsFromFile(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		reader := bytes.NewReader(assets.ExampleJSONLog())
		is, err := source.Reader(reader, config.GetDefaultConfig())
		require.NoError(t, err)
		defer is.Close()
		logEntries, err := is.ParseLogEntries()

		require.NoError(t, err)
		assert.NotEmpty(t, logEntries)
	})

	t.Run("large_line", func(t *testing.T) {
		t.Parallel()

		longLine := strings.Repeat("1", 2*1024*1024)

		reader := strings.NewReader(longLine)
		is, err := source.Reader(reader, config.GetDefaultConfig())
		require.NoError(t, err)
		defer is.Close()
		logEntries, err := is.ParseLogEntries()

		require.NoError(t, err)
		assert.NotEmpty(t, logEntries)
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

	require.Len(t, logEntries.Entries, 0)
}
