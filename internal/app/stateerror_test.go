package app_test

import (
	"fmt"
	"testing"

	"github.com/hedhyw/json-log-viewer/assets"
	"github.com/hedhyw/json-log-viewer/internal/app"
	"github.com/hedhyw/json-log-viewer/internal/pkg/events"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestStateError(t *testing.T) {
	t.Parallel()

	errTest := getTestError()

	setup := func() tea.Model {
		model := newTestModel(t, assets.ExampleJSONLog())
		model = handleUpdate(model, events.ErrorOccuredMsg{Err: errTest})

		_, ok := model.(app.StateErrorModel)
		assert.Truef(t, ok, "%s", model)

		return model
	}

	t.Run("rendered", func(t *testing.T) {
		t.Parallel()

		model := setup()
		rendered := model.View()
		assert.Contains(t, rendered, errTest.Error())
	})

	t.Run("any_key_msg", func(t *testing.T) {
		t.Parallel()

		model := setup()

		_, cmd := model.Update(tea.KeyMsg{})
		assert.Equal(t, tea.Quit(), cmd())
	})

	t.Run("unknown_message", func(t *testing.T) {
		t.Parallel()

		model := setup()

		model, _ = model.Update(nil)

		_, ok := model.(app.StateErrorModel)
		assert.Truef(t, ok, "%s", model)
	})

	t.Run("stringer", func(t *testing.T) {
		t.Parallel()

		model := setup()

		stringer, ok := model.(fmt.Stringer)
		if assert.True(t, ok) {
			assert.Contains(t, stringer.String(), "StateError")
		}
	})
}
