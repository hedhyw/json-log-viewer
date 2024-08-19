package app_test

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/hedhyw/json-log-viewer/assets"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestAppViewResized(t *testing.T) {
	t.Parallel()

	model, source := newTestModel(t, assets.ExampleJSONLog())
	defer source.Close()

	windowSize := tea.WindowSizeMsg{
		Width:  60,
		Height: 10,
	}

	model = handleUpdate(model, windowSize)

	rendered := model.View()
	lines := strings.Split(rendered, "\n")
	if assert.NotEmpty(t, lines, rendered) {
		assert.Less(t, utf8.RuneCountInString(lines[0]), windowSize.Width, rendered)
	}
}
