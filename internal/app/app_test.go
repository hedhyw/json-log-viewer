package app_test

import (
	"errors"
	"os"
	"strings"
	"testing"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hedhyw/json-log-viewer/assets"
	"github.com/hedhyw/json-log-viewer/internal/app"
	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
	"github.com/hedhyw/json-log-viewer/internal/pkg/tests"
)

func TestAppViewFiltered(t *testing.T) {
	const (
		termIncluded = "included"
		termExcluded = "excluded"
	)

	const jsonFile = `
	{"time":"1970-01-01T00:00:00.00","level":"INFO","message": "` + termIncluded + `"}
	{"time":"1970-01-01T00:00:00.00","level":"INFO","message": "` + termExcluded + `"}
	`

	appModel := newTestModel(t, []byte(jsonFile))

	rendered := appModel.View()
	assert.Contains(t, rendered, termIncluded)
	assert.Contains(t, rendered, termExcluded)

	// Open filter.
	appModel, _ = toAppModel(appModel.Update(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'f'},
	}))
	assert.True(t, appModel.IsFilterShown(), appModel.View())

	// Write term to search by.
	appModel, _ = toAppModel(appModel.Update(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(termIncluded),
	}))

	// Filter.
	appModel, cmd := toAppModel(appModel.Update(tea.KeyMsg{
		Type: tea.KeyEnter,
	}))
	assert.False(t, appModel.IsFilterShown(), appModel.View())

	appModel, _ = toAppModel(appModel.Update(cmd()))

	// Assert.
	if assert.True(t, appModel.IsFiltered()) {
		rendered = appModel.View()
		assert.Contains(t, rendered, termIncluded)
		assert.NotContains(t, rendered, termExcluded)
	}
}

func TestAppViewMainScreen(t *testing.T) {
	const jsonFile = `{"time":"1970-01-01T00:00:00.00","level":"INFO","message": "test"}`

	appModel := newTestModel(t, []byte(jsonFile))

	if assert.True(t, appModel.IsTableShown()) {
		rendered := appModel.View()
		assert.Contains(t, rendered, "info")
		assert.Contains(t, rendered, "1970-01-01T00:00:00.00")
		assert.Contains(t, rendered, "test")
	}
}

func TestAppEnterAndCloseJSONView(t *testing.T) {
	appModel := newTestModel(t, assets.ExampleJSONLog())

	appModel, _ = toAppModel(appModel.Update(tea.KeyMsg{Type: tea.KeyEnter}))
	assert.True(t, appModel.IsJSONShown())

	appModel, _ = toAppModel(appModel.Update(tea.KeyMsg{Type: tea.KeyEnter}))
	assert.False(t, appModel.IsJSONShown())
}

func TestAppQuit(t *testing.T) {
	t.Parallel()

	appModel := newTestModel(t, assets.ExampleJSONLog())

	t.Run("ctrl_and_c", func(t *testing.T) {
		t.Parallel()

		_, cmd := toAppModel(appModel.Update(tea.KeyMsg{Type: tea.KeyCtrlC}))

		if assert.NotNil(t, cmd) {
			assert.Equal(t, tea.Quit(), cmd())
		}
	})

	t.Run("esc", func(t *testing.T) {
		t.Parallel()

		_, cmd := toAppModel(appModel.Update(tea.KeyMsg{Type: tea.KeyEsc}))

		if assert.NotNil(t, cmd) {
			assert.Equal(t, tea.Quit(), cmd())
		}
	})

	t.Run("q", func(t *testing.T) {
		t.Parallel()

		_, cmd := toAppModel(appModel.Update(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'q'},
		}))

		if assert.NotNil(t, cmd) {
			assert.Equal(t, tea.Quit(), cmd())
		}
	})
}

func TestAppViewFiltereRunes(t *testing.T) {
	t.Parallel()

	appModel := newTestModel(t, assets.ExampleJSONLog())

	appModel, _ = toAppModel(appModel.Update(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'f'},
	}))
	assert.True(t, appModel.IsFilterShown(), appModel.View())

	appModel, cmd := toAppModel(appModel.Update(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'q'},
	}))
	assert.NotEqual(t, tea.Quit(), cmd())

	appModel, cmd = toAppModel(appModel.Update(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'f'},
	}))
	assert.NotEqual(t, tea.Quit(), cmd())
	assert.True(t, appModel.IsFilterShown(), appModel.View())
}

func TestAppViewReload(t *testing.T) {
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

	appModel := newTestModel(t, []byte(jsonFile))

	rendered := appModel.View()
	assert.NotContains(t, rendered, expected)

	err := os.WriteFile(appModel.File(), []byte(jsonFileUpdated), os.ModePerm)
	require.NoError(t, err)

	appModel, _ = toAppModel(appModel.Update(tea.KeyMsg{
		Type: tea.KeyUp,
	}))

	rendered = appModel.View()
	assert.NotContains(t, rendered, expected)
}

func TestAppViewResized(t *testing.T) {
	t.Parallel()

	appModel := newTestModel(t, assets.ExampleJSONLog())

	windowSize := tea.WindowSizeMsg{
		Width:  60,
		Height: 10,
	}

	appModel, _ = toAppModel(appModel.Update(windowSize))

	rendered := appModel.View()
	lines := strings.Split(rendered, "\n")
	if assert.NotEmpty(t, lines, rendered) {
		assert.Less(t, utf8.RuneCountInString(lines[0]), windowSize.Width, rendered)
	}
}

func TestAppViewError(t *testing.T) {
	t.Parallel()

	appModel := newTestModel(t, assets.ExampleJSONLog())

	// nolint: goerr113 // It is a test.
	errMsg := errors.New("error description")

	appModel, _ = toAppModel(appModel.Update(errMsg))
	assert.True(t, appModel.IsErrorShown())

	rendered := appModel.View()
	assert.Contains(t, rendered, errMsg.Error())
}

func newTestModel(tb testing.TB, content []byte) app.Model {
	tb.Helper()

	testFile := tests.RequireCreateFile(tb, content)

	appModel := app.NewModel(testFile, config.GetDefaultConfig())
	cmd := appModel.Init()

	appModel, _ = toAppModel(appModel.Update(cmd()))

	return appModel
}

func toAppModel(teaModel tea.Model, cmd tea.Cmd) (app.Model, tea.Cmd) {
	appModel, _ := teaModel.(app.Model)

	return appModel, cmd
}

func TestAppViewFilterClear(t *testing.T) {
	const termIncluded = "included"

	const jsonFile = `
	{"time":"1970-01-01T00:00:00.00","level":"INFO","message": "` + termIncluded + `"}
	`

	appModel := newTestModel(t, []byte(jsonFile))

	rendered := appModel.View()
	assert.Contains(t, rendered, termIncluded)

	// Open filter.
	appModel, _ = toAppModel(appModel.Update(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'f'},
	}))
	assert.True(t, appModel.IsFilterShown(), appModel.View())

	// Filter to exclude everything.
	appModel, _ = toAppModel(appModel.Update(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(termIncluded + "_not_found"),
	}))
	appModel, cmd := toAppModel(appModel.Update(tea.KeyMsg{
		Type: tea.KeyEnter,
	}))
	assert.False(t, appModel.IsFilterShown(), appModel.View())

	appModel, _ = toAppModel(appModel.Update(cmd()))

	rendered = appModel.View()
	assert.NotContains(t, rendered, termIncluded)

	// Come back
	appModel, cmd = toAppModel(appModel.Update(tea.KeyMsg{
		Type: tea.KeyEsc,
	}))
	assert.False(t, appModel.IsFilterShown(), appModel.View())

	appModel, _ = toAppModel(appModel.Update(cmd()))

	// Assert.
	if assert.False(t, appModel.IsFiltered()) {
		rendered = appModel.View()
		assert.Contains(t, rendered, termIncluded)
	}
}
