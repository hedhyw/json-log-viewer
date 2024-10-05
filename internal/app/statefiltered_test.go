package app_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hedhyw/json-log-viewer/internal/app"
	"github.com/hedhyw/json-log-viewer/internal/pkg/events"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestStateFiltered(t *testing.T) {
	t.Parallel()

	const (
		termIncluded = "included"
		termExcluded = "excluded"
	)

	const jsonFile = `
	{"time":"1970-01-01T00:00:00.00","level":"INFO","message": "` + termIncluded + `"}
	{"time":"1970-01-01T00:00:00.00","level":"INFO","message": "` + termExcluded + `"}
	`

	setup := func() tea.Model {
		model := newTestModel(t, []byte(jsonFile))

		rendered := model.View()
		assert.Contains(t, rendered, termIncluded)
		assert.Contains(t, rendered, termExcluded)

		// Open filtering.
		model = handleUpdate(model, tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'f'},
		})

		lines := strings.Split(model.View(), "\n")
		assert.Contains(t, lines[len(lines)-1], ">")

		_, ok := model.(app.StateFilteringModel)
		assert.Truef(t, ok, "%s", model)

		// Write term to search by.
		model = handleUpdate(model, tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune(termIncluded[1:]),
		})

		// Filter.
		model = handleUpdate(model, tea.KeyMsg{
			Type: tea.KeyEnter,
		})

		_, ok = model.(app.StateFilteredModel)
		if assert.Truef(t, ok, "%s", model) {
			rendered = model.View()
			assert.Contains(t, rendered, termIncluded)
			assert.Contains(t, rendered, "filtered 1 by: "+termIncluded[1:])
			assert.NotContains(t, rendered, termExcluded)
		}

		return model
	}

	t.Run("reopen_filter", func(t *testing.T) {
		t.Parallel()

		model := setup()
		model = handleUpdate(model, tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'f'},
		})

		_, ok := model.(app.StateFilteringModel)
		assert.Truef(t, ok, "%s", model)
	})

	t.Run("open_hide_json_view", func(t *testing.T) {
		t.Parallel()

		model := setup()
		model = handleUpdate(model, tea.KeyMsg{
			Type: tea.KeyEnter,
		})

		_, ok := model.(app.StateViewRowModel)
		assert.Truef(t, ok, "%s", model)

		model = handleUpdate(model, tea.KeyMsg{
			Type: tea.KeyEsc,
		})

		_, ok = model.(app.StateFilteredModel)
		assert.Truef(t, ok, "%s", model)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		model := setup()
		model = handleUpdate(model, events.ErrorOccuredMsg{Err: getTestError()})

		_, ok := model.(app.StateErrorModel)
		assert.Truef(t, ok, "%s", model)
	})

	t.Run("navigation", func(t *testing.T) {
		t.Parallel()

		model := setup()
		model = handleUpdate(model, tea.KeyMsg{
			Type: tea.KeyUp,
		})

		_, ok := model.(app.StateFilteredModel)
		assert.Truef(t, ok, "%s", model)
	})

	t.Run("returned", func(t *testing.T) {
		t.Parallel()

		model := setup()
		model = handleUpdate(model, tea.KeyMsg{
			Type: tea.KeyEsc,
		})

		_, ok := model.(app.StateLoadedModel)
		assert.Truef(t, ok, "%s", model)
	})

	t.Run("stringer", func(t *testing.T) {
		t.Parallel()

		model := setup()
		stringer, ok := model.(fmt.Stringer)
		if assert.True(t, ok) {
			assert.Contains(t, stringer.String(), "StateFiltered")
		}
	})

	t.Run("updated", func(t *testing.T) {
		t.Parallel()

		model := setup()
		model = handleUpdate(model, &events.LogEntriesUpdateMsg{})

		rendered := model.View()
		assert.Contains(t, rendered, termIncluded)
	})
}
