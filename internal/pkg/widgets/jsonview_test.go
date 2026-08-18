package widgets_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/keymap"
	"github.com/hedhyw/json-log-viewer/internal/pkg/widgets"

	"github.com/stretchr/testify/assert"
)

func TestNewJSONViewModel(t *testing.T) {
	t.Parallel()

	t.Run("plain_text", func(t *testing.T) {
		t.Parallel()

		model, _ := widgets.NewJSONViewModel(
			[]byte(text),
			getFakeTeaWindowSizeMsg(),
			keymap.GetDefaultKeys(),
		)

		_, ok := model.(widgets.PlainLogModel)
		assert.Truef(t, ok, "actual type: %T", model)
	})

	t.Run("json", func(t *testing.T) {
		t.Parallel()

		model, _ := widgets.NewJSONViewModel(
			[]byte(`{"hello":"world"}`),
			getFakeTeaWindowSizeMsg(),
			keymap.GetDefaultKeys(),
		)

		_, ok := model.(widgets.PlainLogModel)
		assert.Falsef(t, ok, "actual type: %T", model)
	})
}

func TestJSONViewModelHelp(t *testing.T) {
	t.Parallel()

	newModel := func(t *testing.T) tea.Model {
		t.Helper()

		model, _ := widgets.NewJSONViewModel(
			[]byte(`{"hello":"world"}`),
			getFakeTeaWindowSizeMsg(),
			keymap.GetDefaultKeys(),
		)

		return model
	}

	// The runes are not shared between the subtests, because a text input
	// of the JSON view sanitizes them in place.
	runeKey := func(r rune) tea.KeyMsg {
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
	}

	t.Run("toggle", func(t *testing.T) {
		t.Parallel()

		model := newModel(t)
		assert.NotContains(t, model.View(), "search regexp")

		model, _ = model.Update(runeKey('?'))
		assert.Contains(t, model.View(), "search regexp")

		// Any key closes the help.
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
		assert.NotContains(t, model.View(), "search regexp")
		assert.Contains(t, model.View(), "hello")
	})

	t.Run("ignored_while_searching", func(t *testing.T) {
		t.Parallel()

		model := newModel(t)

		model, _ = model.Update(runeKey('/'))
		model, _ = model.Update(runeKey('?'))
		assert.NotContains(t, model.View(), "search regexp")

		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEscape})
		model, _ = model.Update(runeKey('?'))
		assert.Contains(t, model.View(), "search regexp")
	})

	t.Run("exit_is_not_intercepted", func(t *testing.T) {
		t.Parallel()

		model := newModel(t)

		model, _ = model.Update(runeKey('?'))

		_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
		if assert.NotNil(t, cmd) {
			assert.Equal(t, tea.Quit(), cmd())
		}
	})
}
