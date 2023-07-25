package source_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hedhyw/json-log-viewer/assets"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
	"github.com/hedhyw/json-log-viewer/internal/pkg/tests"
)

func TestLoadLogsFromFile(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		testFile := tests.RequireCreateFile(t, assets.ExampleJSONLog())

		msg := source.LoadLogsFromFile(testFile)()

		logEntries, ok := msg.(source.LogEntries)
		if assert.Truef(t, ok, "actual type: %T", msg) {
			assert.NotEmpty(t, logEntries)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()

		msg := source.LoadLogsFromFile("not_found_for_" + t.Name())()

		_, ok := msg.(error)
		assert.Truef(t, ok, "actual type: %T", msg)
	})

	t.Run("large_line", func(t *testing.T) {
		t.Parallel()

		longLine := strings.Repeat("1", 2*1024*1024)
		testFile := tests.RequireCreateFile(t, []byte(longLine))

		msg := source.LoadLogsFromFile(testFile)()

		logEntries, ok := msg.(source.LogEntries)
		if assert.Truef(t, ok, "actual type: %T", msg) {
			assert.NotEmpty(t, logEntries)
		}
	})
}
