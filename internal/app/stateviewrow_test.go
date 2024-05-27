package app_test

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hedhyw/json-log-viewer/assets"
	"github.com/hedhyw/json-log-viewer/internal/app"
	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
)

func TestStateViewRow(t *testing.T) {
	t.Parallel()

	model := newTestModel(t, assets.ExampleJSONLog())

	model = handleUpdate(model, tea.KeyMsg{Type: tea.KeyEnter})

	_, ok := model.(app.StateViewRowModel)
	require.Truef(t, ok, "%s", model)

	t.Run("close", func(t *testing.T) {
		t.Parallel()

		model := handleUpdate(model, tea.KeyMsg{Type: tea.KeyEnter})
		_, ok := model.(app.StateLoadedModel)
		require.Truef(t, ok, "%s", model)
	})

	t.Run("stringer", func(t *testing.T) {
		t.Parallel()

		stringer, ok := model.(fmt.Stringer)
		if assert.True(t, ok) {
			assert.Contains(t, stringer.String(), "StateViewRow")
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		model := handleUpdate(model, events.ErrorOccuredMsg{Err: getTestError()})

		_, ok := model.(app.StateErrorModel)
		assert.Truef(t, ok, "%s", model)
	})

	// nolint: tparallel // antonmedv/fx uses mutable model.
	t.Run("navigation", func(t *testing.T) {
		model = handleUpdate(model, tea.KeyMsg{
			Type: tea.KeyRight,
		})

		_, ok := model.(app.StateViewRowModel)
		assert.Truef(t, ok, "%s", model)
	})
}
