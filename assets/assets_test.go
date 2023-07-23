package assets_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hedhyw/json-log-viewer/assets"
)

func TestExampleJSONLog(t *testing.T) {
	t.Parallel()

	content := assets.ExampleJSONLog()
	if assert.NotEmpty(t, content) {
		assert.Contains(t, string(content), "{")
	}
}
