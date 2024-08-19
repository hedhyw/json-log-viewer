package app_test

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hedhyw/json-log-viewer/assets"
	"github.com/hedhyw/json-log-viewer/internal/app"
	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestStateLoadedEmpty(t *testing.T) {
	t.Parallel()

	model, source := newTestModel(t, []byte(""))
	defer source.Close()

	_, ok := model.(app.StateLoadedModel)
	require.Truef(t, ok, "%s", model)

	model, cmd := model.Update(events.EscKeyClicked())
	require.NotNil(t, model)
	requireCmdMsg(t, tea.Quit(), cmd)
}

func TestStateLoaded(t *testing.T) {
	t.Parallel()

	setup := func() (tea.Model, *source.Source) {
		const jsonFile = `{"time":"1970-01-01T00:00:00.00","level":"INFO","message": "test"}`

		model, source := newTestModel(t, []byte(jsonFile))

		_, ok := model.(app.StateLoadedModel)
		require.Truef(t, ok, "%s", model)

		return model, source
	}

	t.Run("stringer", func(t *testing.T) {
		t.Parallel()

		model, source := setup()
		defer source.Close()

		stringer, ok := model.(fmt.Stringer)
		if assert.True(t, ok) {
			assert.Contains(t, stringer.String(), "StateLoaded")
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		model, source := setup()
		defer source.Close()

		model = handleUpdate(model, events.ErrorOccuredMsg{Err: getTestError()})

		_, ok := model.(app.StateErrorModel)
		assert.Truef(t, ok, "%s", model)
	})

	t.Run("version_printed", func(t *testing.T) {
		t.Parallel()
		model, source := setup()
		defer source.Close()

		model = handleUpdate(model, events.HelpKeyClicked())
		view := model.View()
		assert.Contains(t, view, testVersion)
	})
}

func TestStateLoadedQuit(t *testing.T) {
	t.Parallel()

	model, source := newTestModel(t, assets.ExampleJSONLog())
	defer source.Close()

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

/*
go test -benchmem -run=^$ -bench ^BenchmarkStateLoadedBig$ github.com/hedhyw/json-log-viewer/internal/app

goos: linux
goarch: amd64
pkg: github.com/hedhyw/json-log-viewer/internal/app
cpu: 12th Gen Intel(R) Core(TM) i7-1255U
BenchmarkStateLoadedBig-12    	16499398	        78.08 ns/op	     199 B/op	       0 allocs/op.
*/
func BenchmarkStateLoadedBig(b *testing.B) {
	content := strings.Repeat(`{"time":"1970-01-01T00:00:00.00","level":"INFO","message": "test2"}`+"\n", b.N)
	contentReader := strings.NewReader(content)

	cfg := config.GetDefaultConfig()

	model, modelSource := newTestModel(b, []byte(`{}`))
	defer modelSource.Close()

	_, ok := model.(app.StateLoadedModel)
	if !ok {
		b.Fatal(model.View())
	}

	b.ResetTimer()

	is, err := source.Reader(contentReader, cfg)
	require.NoError(b, err)
	defer is.Close()

	logEntries, err := is.ParseLogEntries()
	if err != nil {
		b.Fatal(model.View())
	}

	model.Update(events.LogEntriesUpdateMsg(logEntries))
}

func overwriteFileInStateLoaded(tb testing.TB, model tea.Model, content []byte) {
	tb.Helper()

	stateLoaded, ok := model.(app.StateLoadedModel)
	require.True(tb, ok)

	// nolint: gosec // Test.
	err := os.WriteFile(
		stateLoaded.Application.FileName,
		content,
		os.ModePerm,
	)
	require.NoError(tb, err)
}
