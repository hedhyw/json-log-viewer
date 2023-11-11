package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
)

// StateInitial is an initial loading state.
type StateInitial struct {
	helper
}

func newStateInitial(application Application) StateInitial {
	return StateInitial{
		helper: helper{Application: application},
	}
}

// Init initializes component. It implements tea.Model.
func (s StateInitial) Init() tea.Cmd {
	return s.helper.LoadEntries
}

// View renders component. It implements tea.Model.
func (s StateInitial) View() string {
	return "Loading..."
}

// Update handles events. It implements tea.Model.
func (s StateInitial) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	s.helper = s.helper.Update(msg)

	switch msg := msg.(type) {
	case events.ErrorOccuredMsg:
		return s.handleErrorOccuredMsg(msg)
	case events.LogEntriesLoadedMsg:
		return s.handleLogEntriesLoadedMsg(msg)
	case tea.KeyMsg:
		return s, tea.Quit
	default:
		return s, nil
	}
}

// String implements fmt.Stringer.
func (s StateInitial) String() string {
	return modelValue(s)
}
