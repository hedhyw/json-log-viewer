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

		logEntries, err := source.ParseLogEntriesFromReader(
			bytes.NewReader(assets.ExampleJSONLog()),
			config.GetDefaultConfig(),
		)
		require.NoError(t, err)
		assert.NotEmpty(t, logEntries)
	})

	t.Run("large_line", func(t *testing.T) {
		t.Parallel()

		longLine := strings.Repeat("1", 2*1024*1024)

		logEntries, err := source.ParseLogEntriesFromReader(
			strings.NewReader(longLine),
			config.GetDefaultConfig(),
		)
		require.NoError(t, err)
		assert.NotEmpty(t, logEntries)
	})
}

func TestParseLogEntriesFromReaderLimited(t *testing.T) {
	t.Parallel()

	content := `{}`

	cfg := config.GetDefaultConfig()
	cfg.MaxFileSizeBytes = 1

	logEntries, err := source.ParseLogEntriesFromReader(strings.NewReader(content), cfg)
	require.NoError(t, err)

	if assert.Len(t, logEntries, 1) {
		assert.Equal(t, content[:cfg.MaxFileSizeBytes], string(logEntries[0].Line))
	}
}
