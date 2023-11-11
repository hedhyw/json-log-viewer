package app_test

import (
	"fmt"
	"testing"

	"github.com/hedhyw/json-log-viewer/internal/app"
	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
	"github.com/hedhyw/json-log-viewer/internal/pkg/events"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStateInitial(t *testing.T) {
	t.Parallel()

	model := app.NewModel("", config.GetDefaultConfig())

	_, ok := model.(app.StateInitial)
	require.Truef(t, ok, "%s", model)

	t.Run("stringer", func(t *testing.T) {
		t.Parallel()

		stringer, ok := model.(fmt.Stringer)
		if assert.True(t, ok) {
			assert.Contains(t, stringer.String(), "StateInitial")
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

		_, cmd := model.Update(tea.KeyMsg{
			Type: tea.KeyUp,
		})

		assert.Equal(t, tea.Quit(), cmd())
	})

	t.Run("unknown_update", func(t *testing.T) {
		t.Parallel()

		model := handleUpdate(model, events.ViewRowsReloadRequestedMsg{})

		_, ok := model.(app.StateInitial)
		require.Truef(t, ok, "%s", model)
	})
}
