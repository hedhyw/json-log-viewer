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
	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

func TestStateLoadedEmpty(t *testing.T) {
	t.Parallel()

	model := newTestModel(t, []byte(""))

	_, ok := model.(app.StateLoadedModel)
	require.Truef(t, ok, "%s", model)

	model, cmd := model.Update(events.EscKeyClicked())
	require.NotNil(t, model)
	requireCmdMsg(t, tea.Quit(), cmd)
}

func TestStateLoaded(t *testing.T) {
	t.Parallel()

	setup := func() tea.Model {
		const jsonFile = `{"time":"1970-01-01T00:00:00.00","level":"INFO","message": "test"}`

		model := newTestModel(t, []byte(jsonFile))

		_, ok := model.(app.StateLoadedModel)
		require.Truef(t, ok, "%s", model)

		return model
	}

	t.Run("stringer", func(t *testing.T) {
		t.Parallel()
		model := setup()

		stringer, ok := model.(fmt.Stringer)
		if assert.True(t, ok) {
			assert.Contains(t, stringer.String(), "StateLoaded")
		}
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		model := setup()

		model = handleUpdate(model, events.ErrorOccuredMsg{Err: getTestError()})

		_, ok := model.(app.StateErrorModel)
		assert.Truef(t, ok, "%s", model)
	})

	t.Run("version_printed", func(t *testing.T) {
		t.Parallel()
		model := setup()

		model = handleUpdate(model, events.HelpKeyClicked())
		view := model.View()
		assert.Contains(t, view, testVersion)
	})

	t.Run("hide_help", func(t *testing.T) {
		t.Parallel()
		model := setup()

		model = handleUpdate(model, events.HelpKeyClicked())
		model = handleUpdate(model, events.HelpKeyClicked())

		view := model.View()
		assert.NotContains(t, view, testVersion)
	})

	t.Run("label_following_default", func(t *testing.T) {
		t.Parallel()

		model := setup()

		view := model.View()
		assert.Contains(t, view, "following")
	})

	t.Run("label_not_following", func(t *testing.T) {
		t.Parallel()

		model := setup()
		model = handleUpdate(model, tea.KeyMsg{Type: tea.KeyDown})

		view := model.View()
		assert.NotContains(t, view, "following")
	})

	t.Run("label_reverse_default", func(t *testing.T) {
		t.Parallel()

		model := setup()

		view := model.View()
		assert.Contains(t, view, "reverse")
	})

	t.Run("label_not_reverse", func(t *testing.T) {
		t.Parallel()

		model := setup()
		model = handleUpdate(model, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

		view := model.View()
		assert.NotContains(t, view, "reverse")
	})

	t.Run("label_not_reverse_not_following", func(t *testing.T) {
		t.Parallel()

		model := setup()
		model = handleUpdate(model, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
		model = handleUpdate(model, tea.KeyMsg{Type: tea.KeyDown})

		view := model.View()
		assert.NotContains(t, view, "reverse")
		assert.NotContains(t, view, "following")
	})
}

func TestStateLoadedQuit(t *testing.T) {
	t.Parallel()

	t.Run("ctrl_and_c", func(t *testing.T) {
		t.Parallel()
		model := newTestModel(t, assets.ExampleJSONLog())

		_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
		requireCmdMsg(t, tea.Quit(), cmd)
	})

	t.Run("esc", func(t *testing.T) {
		t.Parallel()
		model := newTestModel(t, assets.ExampleJSONLog())

		_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
		requireCmdMsg(t, tea.Quit(), cmd)
	})

	t.Run("q", func(t *testing.T) {
		t.Parallel()
		model := newTestModel(t, assets.ExampleJSONLog())

		_, cmd := model.Update(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'q'},
		})
		requireCmdMsg(t, tea.Quit(), cmd)
	})

	t.Run("f10", func(t *testing.T) {
		t.Parallel()
		model := newTestModel(t, assets.ExampleJSONLog())

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

	model := newTestModel(b, []byte(`{}`))

	_, ok := model.(app.StateLoadedModel)
	if !ok {
		b.Fatal(model.View())
	}

	b.ResetTimer()

	inputSource, err := source.Reader(contentReader, cfg)
	require.NoError(b, err)
	b.Cleanup(func() { _ = inputSource.Close() })

	logEntries, err := inputSource.ParseLogEntries()
	if err != nil {
		b.Fatal(model.View())
	}

	model.Update(events.LogEntriesUpdateMsg(logEntries))
}
