package source_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hedhyw/json-log-viewer/assets"
	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
	"github.com/hedhyw/json-log-viewer/internal/pkg/tests"
)

func TestLoadLogsFromFile(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		testFile := tests.RequireCreateFile(t, assets.ExampleJSONLog())

		logEntries, err := source.LoadLogsFromFile(
			testFile,
			config.GetDefaultConfig(),
		)
		require.NoError(t, err)
		assert.NotEmpty(t, logEntries)
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()

		_, err := source.LoadLogsFromFile(
			"not_found_for_"+t.Name(),
			config.GetDefaultConfig(),
		)
		assert.Error(t, err)
	})

	t.Run("large_line", func(t *testing.T) {
		t.Parallel()

		longLine := strings.Repeat("1", 2*1024*1024)
		testFile := tests.RequireCreateFile(t, []byte(longLine))

		logEntries, err := source.LoadLogsFromFile(
			testFile,
			config.GetDefaultConfig(),
		)
		require.NoError(t, err)
		assert.NotEmpty(t, logEntries)
	})
}
