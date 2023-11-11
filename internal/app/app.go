package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
)

// Application global state.
type Application struct {
	Path   string
	Config *config.Config

	BaseStyle   lipgloss.Style
	FooterStyle lipgloss.Style

	LastWindowSize tea.WindowSizeMsg
}

func newApplication(path string, config *config.Config) Application {
	const (
		initialWidth  = 70
		initialHeight = 20
	)

	return Application{
		Path:   path,
		Config: config,

		BaseStyle:   getBaseStyle(),
		FooterStyle: getFooterStyle(),

		LastWindowSize: tea.WindowSizeMsg{
			Width:  initialWidth,
			Height: initialHeight,
		},
	}
}

// NewModel initializes a new application model. It accept the path
// to the file with logs.
func NewModel(path string, config *config.Config) tea.Model {
	return newStateInitial(newApplication(path, config))
}
