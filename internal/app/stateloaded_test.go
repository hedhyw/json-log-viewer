package app_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hedhyw/json-log-viewer/assets"
	"github.com/hedhyw/json-log-viewer/internal/app"
	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"

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

	rendered := model.View()
	assert.NotContains(t, rendered, expected)

	overwriteFileInStateLoaded(t, model, []byte(jsonFileUpdated))

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

	t.Run("threshold", func(t *testing.T) {
		t.Parallel()

		model := newTestModel(t, []byte(jsonFile))

		model = handleUpdate(model, tea.KeyMsg{
			Type: tea.KeyUp,
		})

		overwriteFileInStateLoaded(t, model, []byte(jsonFileUpdated))

		model = handleUpdate(model, tea.KeyMsg{
			Type: tea.KeyUp,
		})

		rendered := model.View()
		assert.NotContains(t, rendered, expected)
	})
}

/*
go test -benchmem -run=^$ -bench ^BenchmarkStateLoadedBig$ github.com/hedhyw/json-log-viewer/internal/app

goos: linux
goarch: amd64
pkg: github.com/hedhyw/json-log-viewer/internal/app
cpu: 12th Gen Intel(R) Core(TM) i7-1255U
BenchmarkStateLoadedBig-12    	16499398	        78.08 ns/op	     199 B/op	       0 allocs/op
*/
func BenchmarkStateLoadedBig(b *testing.B) {
	content := strings.Repeat(`{"time":"1970-01-01T00:00:00.00","level":"INFO","message": "test2"}`+"\n", b.N)
	contentReader := strings.NewReader(content)

	cfg := config.GetDefaultConfig()

	model := newTestModel(b, []byte(`{}`))

	_, ok := model.(app.StateLoaded)
	if !ok {
		b.Fatal(model.View())
	}

	b.ResetTimer()

	logEntries, err := source.ParseLogEntriesFromReader(contentReader, cfg)
	if err != nil {
		b.Fatal(model.View())
	}

	model.Update(events.LogEntriesLoadedMsg(logEntries))
}

func overwriteFileInStateLoaded(tb testing.TB, model tea.Model, content []byte) {
	tb.Helper()

	stateLoaded, ok := model.(app.StateLoaded)
	require.True(tb, ok)

	err := os.WriteFile(
		stateLoaded.Application().Path,
		content,
		os.ModePerm,
	)
	require.NoError(tb, err)
}
