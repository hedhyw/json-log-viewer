package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
	"github.com/hedhyw/json-log-viewer/internal/pkg/widgets"
)

// StateViewRowModel is a state that shows extended JSON view.
type StateViewRowModel struct {
	*Application

	previousState stateModel
	initCmd       tea.Cmd

	logEntry source.LogEntry
	jsonView tea.Model

	keys KeyMap
}

func newStateViewRow(
	logEntry source.LogEntry,
	previousState stateModel,
) StateViewRowModel {
	jsonViewModel, cmd := widgets.NewJSONViewModel(logEntry.Line, previousState.getApplication().LastWindowSize)

	return StateViewRowModel{
		Application: previousState.getApplication(),

		previousState: previousState,
		initCmd:       cmd,

		logEntry: logEntry,
		jsonView: jsonViewModel,

		keys: defaultKeys,
	}
}

// Init initializes component. It implements tea.Model.
func (s StateViewRowModel) Init() tea.Cmd {
	return s.initCmd
}

// View renders component. It implements tea.Model.
func (s StateViewRowModel) View() string {
	return s.jsonView.View()
}

// Update handles events. It implements tea.Model.
func (s StateViewRowModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	s.Application.Update(msg)

	switch msg := msg.(type) {
	case events.ErrorOccuredMsg:
		return s.handleErrorOccuredMsg(msg)
	case tea.KeyMsg:
		if key.Matches(msg, s.keys.Back) {
			return s.previousState.refresh()
		}
	}

	s.jsonView, cmd = s.jsonView.Update(msg)

	return s, cmd
}

// String implements fmt.Stringer.
func (s StateViewRowModel) String() string {
	return modelValue(s)
}
