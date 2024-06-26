package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

// Application global state.
type Application struct {
	SourceInput source.Input
	Config      *config.Config

	BaseStyle   lipgloss.Style
	FooterStyle lipgloss.Style

	LastWindowSize tea.WindowSizeMsg

	Version string
}

func newApplication(
	sourceInput source.Input,
	config *config.Config,
	version string,
) Application {
	const (
		initialWidth  = 70
		initialHeight = 20
	)

	return Application{
		SourceInput: sourceInput,
		Config:      config,

		BaseStyle:   getBaseStyle(),
		FooterStyle: getFooterStyle(),

		LastWindowSize: tea.WindowSizeMsg{
			Width:  initialWidth,
			Height: initialHeight,
		},

		Version: version,
	}
}

// NewModel initializes a new application model. It accept the path
// to the file with logs.
func NewModel(
	sourceInput source.Input,
	config *config.Config,
	version string,
) tea.Model {
	return newStateInitial(newApplication(sourceInput, config, version))
}
