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

	setup := func(configSetters ...configSetter) tea.Model {
		const jsonFile = `{"time":"1970-01-01T00:00:00.00","level":"INFO","message": "test"}`

		model := newTestModel(t, []byte(jsonFile), configSetters...)

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

	t.Run("label_reverse_disabled_in_config", func(t *testing.T) {
		t.Parallel()

		model := setup(func(cfg *config.Config) {
			cfg.IsReverseDefault = false
		})

		view := model.View()
		assert.NotContains(t, view, "reverse")
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

func TestStateLoadedNewerEntriesCount(t *testing.T) {
	t.Parallel()

	const entriesCount = 5

	setup := func(configSetters ...configSetter) tea.Model {
		lines := make([]string, 0, entriesCount)

		for i := range entriesCount {
			lines = append(lines, fmt.Sprintf(
				`{"time":"1970-01-01T00:00:00.00","level":"INFO","message":"test %d"}`,
				i,
			))
		}

		model := newTestModel(t, []byte(strings.Join(lines, "\n")), configSetters...)

		_, ok := model.(app.StateLoadedModel)
		require.Truef(t, ok, "%s", model)

		return model
	}

	t.Run("hidden_while_following", func(t *testing.T) {
		t.Parallel()

		model := setup()

		assert.NotContains(t, model.View(), "newer")
	})

	t.Run("counts_entries_after_cursor", func(t *testing.T) {
		t.Parallel()

		model := setup(func(cfg *config.Config) {
			cfg.IsReverseDefault = false
		})

		// The cursor is on the last entry, so there is nothing newer.
		model = handleUpdate(model, tea.KeyMsg{Type: tea.KeyUp})
		assert.Contains(t, model.View(), "1 newer")

		model = handleUpdate(model, tea.KeyMsg{Type: tea.KeyUp})
		assert.Contains(t, model.View(), "2 newer")

		model = handleUpdate(model, tea.KeyMsg{Type: tea.KeyDown})
		assert.Contains(t, model.View(), "1 newer")
	})

	t.Run("hidden_on_the_newest_entry", func(t *testing.T) {
		t.Parallel()

		model := setup(func(cfg *config.Config) {
			cfg.IsReverseDefault = false
		})

		model = handleUpdate(model, tea.KeyMsg{Type: tea.KeyUp})
		model = handleUpdate(model, tea.KeyMsg{Type: tea.KeyDown})

		assert.NotContains(t, model.View(), "newer")
	})

	t.Run("reverse", func(t *testing.T) {
		t.Parallel()

		model := setup()

		// In the reverse mode the newest entry is on the top.
		model = handleUpdate(model, tea.KeyMsg{Type: tea.KeyDown})
		assert.Contains(t, model.View(), "1 newer")
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

	b.Cleanup(func() { assert.NoError(b, inputSource.Close()) })

	logEntries, err := inputSource.ParseLogEntries()
	if err != nil {
		b.Fatal(model.View())
	}

	model.Update(events.LogEntriesUpdateMsg(logEntries))
}
