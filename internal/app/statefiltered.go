package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

// StateFilteredModel is a state that shows filtered records.
type StateFilteredModel struct {
	*Application

	previousState StateLoadedModel
	table         logsTableModel
	logEntries    source.LazyLogEntries

	filterText string
}

func newStateFiltered(
	previousState StateLoadedModel,
	filterText string,
) StateFilteredModel {
	return StateFilteredModel{
		Application: previousState.Application,

		previousState: previousState,

		filterText: filterText,
	}
}

// Init initializes component. It implements tea.Model.
func (s StateFilteredModel) Init() tea.Cmd {
	return func() tea.Msg {
		return &s
	}
}

// View renders component. It implements tea.Model.
func (s StateFilteredModel) View() string {
	footer := s.Application.FooterStyle.Render(
		fmt.Sprintf("filtered %d by: %s", s.logEntries.Len(), s.filterText),
	)

	return s.BaseStyle.Render(s.table.View()) + "\n" + footer
}

// Update handles events. It implements tea.Model.
func (s StateFilteredModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmdBatch []tea.Cmd

	s.Application.Update(msg)

	switch typedMsg := msg.(type) {
	case *StateFilteredModel:
		entries, err := s.Application.Entries.Filter(s.filterText)
		if err != nil {
			return s, events.ShowError(err)
		}
		s.logEntries = entries
		s.table = newLogsTableModel(s.Application, entries)
		msg = events.LogEntriesUpdateMsg(entries)
	case events.LogEntriesUpdateMsg:
		entries, err := s.Application.Entries.Filter(s.filterText)
		if err != nil {
			return s, events.ShowError(err)
		}
		s.logEntries = entries
		msg = events.LogEntriesUpdateMsg(entries)

	case events.ErrorOccuredMsg:
		return s.handleErrorOccuredMsg(typedMsg)
	case events.OpenJSONRowRequestedMsg:
		return s.handleOpenJSONRowRequestedMsg(typedMsg, s)
	case tea.KeyMsg:
		switch {
		case key.Matches(typedMsg, s.keys.Back):
			return s.previousState.refresh()
		case key.Matches(typedMsg, s.keys.Filter):
			return s.handleFilterKeyClickedMsg()
		case key.Matches(typedMsg, s.keys.ToggleViewArrow), key.Matches(typedMsg, s.keys.Open):
			return s.handleRequestOpenJSON()
		}
		if cmd := s.handleKeyMsg(typedMsg); cmd != nil {
			return s, cmd
		}
	default:
		s.table, cmdBatch = batched(s.table.Update(typedMsg))(cmdBatch)
	}

	s.table, cmdBatch = batched(s.table.Update(msg))(cmdBatch)
	return s, tea.Batch(cmdBatch...)
}

func (s StateFilteredModel) handleFilterKeyClickedMsg() (tea.Model, tea.Cmd) {
	state := newStateFiltering(s.previousState)
	return initializeModel(state)
}

func (s StateFilteredModel) handleRequestOpenJSON() (tea.Model, tea.Cmd) {
	if s.logEntries.Len() == 0 {
		return s, events.BackKeyClicked
	}

	return s, events.OpenJSONRowRequested(s.logEntries, s.table.Cursor())
}

func (s StateFilteredModel) getApplication() *Application {
	return s.Application
}

func (s StateFilteredModel) refresh() (stateModel, tea.Cmd) {
	var cmd tea.Cmd
	s.table, cmd = s.table.Update(s.Application.LastWindowSize)
	return s, cmd
}

// String implements fmt.Stringer.
func (s StateFilteredModel) String() string {
	return modelValue(s)
}
