package app_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hedhyw/json-log-viewer/assets"
	"github.com/hedhyw/json-log-viewer/internal/app"
	"github.com/hedhyw/json-log-viewer/internal/pkg/events"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStateLoadedEmpty(t *testing.T) {
	t.Parallel()

	model := newTestModel(t, []byte(""))

	_, ok := model.(app.StateLoaded)
	require.Truef(t, ok, "%s", model)

	model, cmd := model.Update(events.EnterKeyClicked())
	require.NotNil(t, model)
	requireCmdMsg(t, tea.Quit(), cmd)
}

func TestStateLoaded(t *testing.T) {
	t.Parallel()

	const jsonFile = `{"time":"1970-01-01T00:00:00.00","level":"INFO","message": "test"}`

	model := newTestModel(t, []byte(jsonFile))

	_, ok := model.(app.StateLoaded)
	require.Truef(t, ok, "%s", model)

	t.Run("stringer", func(t *testing.T) {
		t.Parallel()

		stringer, ok := model.(fmt.Stringer)
		if assert.True(t, ok) {
			assert.Contains(t, stringer.String(), "StateLoaded")
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		model := handleUpdate(model, events.ErrorOccuredMsg{Err: getTestError()})

		_, ok = model.(app.StateError)
		assert.Truef(t, ok, "%s", model)
	})
}

func TestStateLoadedQuit(t *testing.T) {
	t.Parallel()

	model := newTestModel(t, assets.ExampleJSONLog())

	t.Run("ctrl_and_c", func(t *testing.T) {
		t.Parallel()

		_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
		requireCmdMsg(t, tea.Quit(), cmd)
	})

	t.Run("esc", func(t *testing.T) {
		t.Parallel()

		_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
		requireCmdMsg(t, tea.Quit(), cmd)
	})

	t.Run("q", func(t *testing.T) {
		t.Parallel()

		_, cmd := model.Update(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'q'},
		})
		requireCmdMsg(t, tea.Quit(), cmd)
	})

	t.Run("f10", func(t *testing.T) {
		t.Parallel()

		_, cmd := model.Update(tea.KeyMsg{
			Type: tea.KeyF10,
		})
		requireCmdMsg(t, tea.Quit(), cmd)
	})
}

func TestStateLoadedReload(t *testing.T) {
	t.Parallel()

	const expected = "included"

	const (
		jsonFile = `
		{"time":"1970-01-01T00:00:00.00","level":"INFO","message": "test2"}
		{"time":"1970-01-01T00:00:00.00","level":"INFO","message": "test1"}
		`

		jsonFileUpdated = `
		{"time":"1970-01-01T00:00:00.00","level":"INFO","message": "` + expected + `"}
		` + jsonFile
	)

	model := newTestModel(t, []byte(jsonFile))

	stateLoaded, ok := model.(app.StateLoaded)
	require.True(t, ok)

	rendered := model.View()
	assert.NotContains(t, rendered, expected)

	err := os.WriteFile(
		stateLoaded.Application().Path,
		[]byte(jsonFileUpdated),
		os.ModePerm,
	)
	require.NoError(t, err)

	t.Run("up", func(t *testing.T) {
		t.Parallel()

		model := handleUpdate(model, tea.KeyMsg{
			Type: tea.KeyUp,
		})

		rendered := model.View()
		assert.Contains(t, rendered, expected)
	})

	t.Run("up_down_up_up", func(t *testing.T) {
		t.Parallel()

		// Go from the first row to the second and back.
		model := handleUpdate(model, tea.KeyMsg{
			Type: tea.KeyDown,
		})
		model = handleUpdate(model, tea.KeyMsg{
			Type: tea.KeyUp,
		})
		assert.NotContains(t, rendered, expected)

		// Press Up, there are no rows.
		model = handleUpdate(model, tea.KeyMsg{
			Type: tea.KeyUp,
		})

		rendered := model.View()
		assert.Contains(t, rendered, expected)
	})
}
