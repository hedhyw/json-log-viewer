package app_test

import (
	"fmt"
	"testing"

	"github.com/hedhyw/json-log-viewer/assets"
	"github.com/hedhyw/json-log-viewer/internal/app"
	"github.com/hedhyw/json-log-viewer/internal/pkg/events"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStateError(t *testing.T) {
	t.Parallel()

	errTest := getTestError()

	model := newTestModel(t, assets.ExampleJSONLog())
	model = handleUpdate(model, events.ErrorOccuredMsg{Err: errTest})

	_, ok := model.(app.StateError)
	assert.Truef(t, ok, "%s", model)

	t.Run("rendered", func(t *testing.T) {
		t.Parallel()

		rendered := model.View()
		assert.Contains(t, rendered, errTest.Error())
	})

	t.Run("any_key_msg", func(t *testing.T) {
		t.Parallel()

		_, cmd := model.Update(tea.KeyMsg{})
		assert.Equal(t, tea.Quit(), cmd())
	})

	t.Run("stringer", func(t *testing.T) {
		t.Parallel()

		stringer, ok := model.(fmt.Stringer)
		if assert.True(t, ok) {
			assert.Contains(t, stringer.String(), "StateError")
		}
	})

	t.Run("unknown_update", func(t *testing.T) {
		t.Parallel()

		model := handleUpdate(model, events.ViewRowsReloadRequestedMsg{})

		_, ok := model.(app.StateError)
		require.Truef(t, ok, "%s", model)
	})
}
