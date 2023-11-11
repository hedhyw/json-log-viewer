package app_test

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hedhyw/json-log-viewer/assets"
	"github.com/hedhyw/json-log-viewer/internal/app"
	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
)

func TestStateFiltering(t *testing.T) {
	t.Parallel()

	model := newTestModel(t, assets.ExampleJSONLog())

	model = handleUpdate(model, tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'f'},
	})
	_, ok := model.(app.StateFiltering)
	assert.Truef(t, ok, "%s", model)

	t.Run("input_hotkeys", func(t *testing.T) {
		t.Parallel()

		model := handleUpdate(model, tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'q'},
		})

		model = handleUpdate(model, tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'f'},
		})

		_, ok := model.(app.StateFiltering)
		assert.Truef(t, ok, "%s", model)
	})

	t.Run("returned", func(t *testing.T) {
		t.Parallel()

		model := handleUpdate(model, tea.KeyMsg{
			Type: tea.KeyEsc,
		})

		_, ok := model.(app.StateLoaded)
		assert.Truef(t, ok, "%s", model)
	})

	t.Run("empty_input", func(t *testing.T) {
		t.Parallel()

		model := handleUpdate(model, tea.KeyMsg{
			Type: tea.KeyEnter,
		})

		_, ok := model.(app.StateLoaded)
		require.Truef(t, ok, "%s", model)
	})

	t.Run("stringer", func(t *testing.T) {
		t.Parallel()

		stringer, ok := model.(fmt.Stringer)
		if assert.True(t, ok) {
			assert.Contains(t, stringer.String(), "StateFiltering")
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		model := handleUpdate(model, events.ErrorOccuredMsg{Err: getTestError()})

		_, ok := model.(app.StateError)
		assert.Truef(t, ok, "%s", model)
	})

	t.Run("navigation", func(t *testing.T) {
		t.Parallel()

		model := handleUpdate(model, tea.KeyMsg{
			Type: tea.KeyUp,
		})

		_, ok := model.(app.StateFiltering)
		assert.Truef(t, ok, "%s", model)
	})
}

func TestStateFilteringReset(t *testing.T) {
	t.Parallel()

	const termIncluded = "included"

	const jsonFile = `
	{"time":"1970-01-01T00:00:00.00","level":"INFO","message": "` + termIncluded + `"}
	`

	model := newTestModel(t, []byte(jsonFile))

	rendered := model.View()
	assert.Contains(t, rendered, termIncluded)

	// Open filter.
	model = handleUpdate(model, tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'f'},
	})

	_, ok := model.(app.StateFiltering)
	assert.Truef(t, ok, "%s", model)

	// Filter to exclude everything.
	model = handleUpdate(model, tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(termIncluded + "_not_found"),
	})
	model = handleUpdate(model, tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	_, ok = model.(app.StateFiltered)
	assert.Truef(t, ok, "%s", model)

	t.Run("record_not_included", func(t *testing.T) {
		t.Parallel()

		rendered := model.View()

		index := strings.Index(rendered, "filtered by:")
		if assert.Greater(t, index, 0) {
			rendered = rendered[:index]
		}

		assert.NotContains(t, rendered, termIncluded)

		// Come back
		model := handleUpdate(model, tea.KeyMsg{
			Type: tea.KeyEsc,
		})

		_, ok = model.(app.StateLoaded)
		assert.Truef(t, ok, "%s", model)

		// Assert.
		rendered = model.View()
		assert.Contains(t, rendered, termIncluded)
	})

	t.Run("record_not_included", func(t *testing.T) {
		t.Parallel()

		// Try to open a record where there are no records.
		model := handleUpdate(model, tea.KeyMsg{
			Type: tea.KeyEnter,
		})

		assert.NotNil(t, model)

		_, ok := model.(app.StateLoaded)
		assert.Truef(t, ok, "%s", model)
	})
}
