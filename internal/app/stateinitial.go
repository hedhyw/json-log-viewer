package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
)

// StateInitialModel is an initial loading state.
type StateInitialModel struct {
	*Application
}

func newStateInitial(application *Application) StateInitialModel {
	return StateInitialModel{
		Application: application,
	}
}

// Init initializes component. It implements tea.Model.
func (s StateInitialModel) Init() tea.Cmd {
	return nil
}

// View renders component. It implements tea.Model.
func (s StateInitialModel) View() string {
	return "Loading..."
}

// Update handles events. It implements tea.Model.
func (s StateInitialModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	s.Application.Update(msg)

	switch msg := msg.(type) {
	case events.ErrorOccuredMsg:
		return s.handleErrorOccuredMsg(msg)
	case events.LogEntriesUpdateMsg:
		return s.handleInitialLogEntriesLoadedMsg(msg)
	case tea.KeyMsg:
		return s, tea.Quit
	default:
		return s, nil
	}
}

// String implements fmt.Stringer.
func (s StateInitialModel) String() string {
	return modelValue(s)
}
