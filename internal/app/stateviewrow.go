package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
	"github.com/hedhyw/json-log-viewer/internal/pkg/widgets"
)

// StateViewRow is a state that shows extended JSON view.
type StateViewRow struct {
	helper

	previousState state
	initCmd       tea.Cmd

	logEntry source.LogEntry
	jsonView tea.Model
}

func newStateViewRow(
	application Application,
	logEntry source.LogEntry,
	previousState state,
) StateViewRow {
	jsonViewModel, cmd := widgets.NewJSONViewModel(logEntry.Line, application.LastWindowSize)

	return StateViewRow{
		helper: helper{Application: application},

		previousState: previousState,
		initCmd:       cmd,

		logEntry: logEntry,
		jsonView: jsonViewModel,
	}
}

// Init initializes component. It implements tea.Model.
func (s StateViewRow) Init() tea.Cmd {
	return s.initCmd
}

// View renders component. It implements tea.Model.
func (s StateViewRow) View() string {
	return s.jsonView.View()
}

// Update handles events. It implements tea.Model.
func (s StateViewRow) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	s.helper = s.helper.Update(msg)

	switch msg := msg.(type) {
	case events.ErrorOccuredMsg:
		return s.handleErrorOccuredMsg(msg)
	case events.BackKeyClickedMsg:
		return s.previousState.withApplication(s.Application)
	case events.EnterKeyClickedMsg:
		return s.previousState.withApplication(s.Application)
	case tea.KeyMsg:
		if cmd = s.handleKeyMsg(msg); cmd != nil {
			return s, cmd
		}
	}

	s.jsonView, cmd = s.jsonView.Update(msg)

	return s, cmd
}

func (s StateViewRow) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	if s.isArrowRightKeyMap(msg) {
		return nil
	}

	return s.helper.handleKeyMsg(msg)
}

// String implements fmt.Stringer.
func (s StateViewRow) String() string {
	return modelValue(s)
}
