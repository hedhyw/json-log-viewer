package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

// StateFiltered is a state that shows filtered records.
type StateFiltered struct {
	helper

	previousState StateLoaded
	table         logsTableModel
	logEntries    source.LazyLogEntries

	filterText string
	keys       KeyMap
}

func newStateFiltered(
	application Application,
	previousState StateLoaded,
	filterText string,
) StateFiltered {
	return StateFiltered{
		helper: helper{Application: application},

		previousState: previousState,
		table:         previousState.table,

		filterText: filterText,
		keys:       defaultKeys,
	}
}

// Init initializes component. It implements tea.Model.
func (s StateFiltered) Init() tea.Cmd {
	return func() tea.Msg {
		return events.LogEntriesLoadedMsg(
			s.previousState.logEntries.Filter(s.filterText),
		)
	}
}

// View renders component. It implements tea.Model.
func (s StateFiltered) View() string {
	footer := s.Application.FooterStyle.Render(" filtered by: " + s.filterText)

	return s.BaseStyle.Render(s.table.View()) + "\n" + footer
}

// Update handles events. It implements tea.Model.
func (s StateFiltered) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmdBatch []tea.Cmd

	s.helper = s.helper.Update(msg)

	switch msg := msg.(type) {
	case events.ErrorOccuredMsg:
		return s.handleErrorOccuredMsg(msg)
	case events.LogEntriesLoadedMsg:
		return s.handleLogEntriesLoadedMsg(msg)
	case events.OpenJSONRowRequestedMsg:
		return s.handleOpenJSONRowRequestedMsg(msg, s)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.Back):
			return s.previousState.withApplication(s.Application)
		case key.Matches(msg, s.keys.Filter):
			return s.handleFilterKeyClickedMsg()
		case key.Matches(msg, s.keys.ToggleViewArrow), key.Matches(msg, s.keys.ToggleView):
			return s.handleRequestOpenJSON()
		}
		if cmd := s.handleKeyMsg(msg); cmd != nil {
			return s, cmd
		}
	default:
		s.table, cmdBatch = batched(s.table.Update(msg))(cmdBatch)
	}

	s.table, cmdBatch = batched(s.table.Update(msg))(cmdBatch)

	return s, tea.Batch(cmdBatch...)
}

func (s StateFiltered) handleLogEntriesLoadedMsg(
	msg events.LogEntriesLoadedMsg,
) (tea.Model, tea.Cmd) {
	s.logEntries = source.LazyLogEntries(msg)
	s.table = newLogsTableModel(s.Application, s.logEntries)

	return s, s.table.Init()
}

func (s StateFiltered) handleFilterKeyClickedMsg() (tea.Model, tea.Cmd) {
	state := newStateFiltering(
		s.Application,
		s.previousState,
	)

	return initializeModel(state)
}

func (s StateFiltered) handleRequestOpenJSON() (tea.Model, tea.Cmd) {
	if len(s.logEntries) == 0 {
		return s, events.BackKeyClicked
	}

	return s, events.OpenJSONRowRequested(s.logEntries, s.table.Cursor())
}

func (s StateFiltered) withApplication(application Application) (state, tea.Cmd) {
	s.Application = application

	var cmd tea.Cmd
	s.table, cmd = s.table.Update(s.Application.LastWindowSize)

	return s, cmd
}

// String implements fmt.Stringer.
func (s StateFiltered) String() string {
	return modelValue(s)
}
