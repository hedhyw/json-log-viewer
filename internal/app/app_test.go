package app_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	"github.com/hedhyw/json-log-viewer/assets"
	"github.com/hedhyw/json-log-viewer/internal/app"
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

func newTestModel(tb testing.TB, content []byte) app.Model {
	tb.Helper()

	testFile := tests.RequireCreateFile(tb, content)

	appModel := app.NewModel(testFile)
	cmd := appModel.Init()

	appModel, _ = toAppModel(appModel.Update(cmd()))

	return appModel
}

func toAppModel(teaModel tea.Model, cmd tea.Cmd) (app.Model, tea.Cmd) {
	appModel, _ := teaModel.(app.Model)

	return appModel, cmd
}

func TestAppViewFiltereRunes(t *testing.T) {
	appModel := newTestModel(t, []byte(assets.ExampleJSONLog()))

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
