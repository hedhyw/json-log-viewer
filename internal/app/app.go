package app

import (
	"sync"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hedhyw/json-log-viewer/internal/keymap"
	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"

	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
)

// Application global state.
type Application struct {
	lock *sync.Mutex

	FileName string
	Config   *config.Config

	BaseStyle   lipgloss.Style
	FooterStyle lipgloss.Style

	lastWindowSize tea.WindowSizeMsg
	entries        source.LazyLogEntries
	Version        string

	keys keymap.KeyMap
	help help.Model
}

func newApplication(
	fileName string,
	config *config.Config,
	version string,
) Application {
	const (
		initialWidth  = 70
		initialHeight = 20
	)

	return Application{
		lock: &sync.Mutex{},

		FileName: fileName,
		Config:   config,

		BaseStyle:   getBaseStyle(),
		FooterStyle: getFooterStyle(),

		lastWindowSize: tea.WindowSizeMsg{
			Width:  initialWidth,
			Height: initialHeight,
		},

		Version: version,
		keys:    keymap.GetDefaultKeys(),
		help:    help.New(),
	}
}

// NewModel initializes a new application model. It accept the path
// to the file with logs.
func NewModel(
	fileName string,
	config *config.Config,
	version string,
) tea.Model {
	application := newApplication(fileName, config, version)

	return newStateInitial(&application)
}

func (app *Application) getLogLevelStyle(
	renderedRows []table.Row,
	baseStyle lipgloss.Style,
	rowID int,
) lipgloss.Style {
	if rowID < 0 || rowID >= len(renderedRows) {
		return baseStyle
	}

	row := renderedRows[rowID]

	color := getColorForLogLevel(app.getLogLevelFromLogRow(row))
	if color == "" {
		return baseStyle
	}

	return baseStyle.Foreground(color)
}

// Update application state.
func (app *Application) Update(msg tea.Msg) {
	app.lock.Lock()
	defer app.lock.Unlock()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		app.lastWindowSize = msg
	case events.LogEntriesUpdateMsg:
		app.entries = source.LazyLogEntries(msg)
	}
}

// Entries getter
func (app *Application) Entries() source.LazyLogEntries {
	app.lock.Lock()
	defer app.lock.Unlock()

	return app.entries
}

// LastWindowSize getter
func (app *Application) LastWindowSize() tea.WindowSizeMsg {
	app.lock.Lock()
	defer app.lock.Unlock()

	return app.lastWindowSize
}
