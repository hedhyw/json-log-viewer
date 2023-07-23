package source_test

import (
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
}
