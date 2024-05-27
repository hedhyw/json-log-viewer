package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
)

// StateInitialModel is an initial loading state.
type StateInitialModel struct {
	helper
}

func newStateInitial(application Application) StateInitialModel {
	return StateInitialModel{
		helper: helper{Application: application},
	}
}

// Init initializes component. It implements tea.Model.
func (s StateInitialModel) Init() tea.Cmd {
	return s.helper.LoadEntries
}

// View renders component. It implements tea.Model.
func (s StateInitialModel) View() string {
	return "Loading..."
}

// Update handles events. It implements tea.Model.
func (s StateInitialModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	s.helper = s.helper.Update(msg)

	switch msg := msg.(type) {
	case events.ErrorOccuredMsg:
		return s.handleErrorOccuredMsg(msg)
	case events.LogEntriesLoadedMsg:
		return s.handleLogEntriesLoadedMsg(msg, time.UnixMilli(0))
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
