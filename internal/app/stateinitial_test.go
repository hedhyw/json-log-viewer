package app_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/hedhyw/json-log-viewer/internal/app"
	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStateInitial(t *testing.T) {
	t.Parallel()

	inputSource, err := source.Reader(bytes.NewReader([]byte{}), config.GetDefaultConfig())
	require.NoError(t, err)

	t.Cleanup(func() { assert.NoError(t, inputSource.Close()) })

	model := app.NewModel(
		"-",
		config.GetDefaultConfig(),
		testVersion,
	)

	entries, err := inputSource.ParseLogEntries()
	require.NoError(t, err)
	handleUpdate(model, events.LogEntriesUpdateMsg(entries))

	_, ok := model.(app.StateInitialModel)
	require.Truef(t, ok, "%s", model)

	t.Run("Init", func(t *testing.T) {
		t.Parallel()

		model, ok := model.(app.StateInitialModel)
		require.Truef(t, ok, "%s", model)

		assert.Nil(t, model.Init())
	})

	t.Run("Unknown_Event", func(t *testing.T) {
		t.Parallel()

		model := handleUpdate(model, nil)

		model, ok := model.(app.StateInitialModel)
		require.Truef(t, ok, "%s", model)
	})

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

		_, ok := model.(app.StateErrorModel)
		assert.Truef(t, ok, "%s", model)
	})

	t.Run("navigation", func(t *testing.T) {
		t.Parallel()

		_, cmd := model.Update(tea.KeyMsg{
			Type: tea.KeyUp,
		})

		assert.Equal(t, tea.Quit(), cmd())
	})
}
