package app

import (
	"sync"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hedhyw/json-log-viewer/internal/keymap"
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
