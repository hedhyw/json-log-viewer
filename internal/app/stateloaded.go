package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

const defaultFooter = "[Ctrl+C] Exit; [Esc] Back; [Enter] Open/Hide; [↑↓] Navigation; [F] Filter"

// StateLoaded is a state that shows all loaded records.
type StateLoaded struct {
	helper

	initCmd tea.Cmd

	table      logsTableModel
	logEntries source.LogEntries
}

func newStateViewLogs(application Application, logEntries source.LogEntries) StateLoaded {
	table := newLogsTableModel(application, logEntries)

	return StateLoaded{
		helper: helper{Application: application},

		initCmd: table.Init(),

		table:      table,
		logEntries: logEntries,
	}
}

// Init initializes component. It implements tea.Model.
func (s StateLoaded) Init() tea.Cmd {
	return s.initCmd
}

// View renders component. It implements tea.Model.
func (s StateLoaded) View() string {
	return s.viewTable() + "\n" + s.viewFooter()
}

func (s StateLoaded) viewTable() string {
	return s.BaseStyle.Render(s.table.View())
}

func (s StateLoaded) viewFooter() string {
	return s.FooterStyle.Render(defaultFooter)
}

// Update handles events. It implements tea.Model.
//
// nolint: cyclop // Many events in switch case.
func (s StateLoaded) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmdBatch []tea.Cmd

	s.helper = s.helper.Update(msg)

	switch msg := msg.(type) {
	case events.ErrorOccuredMsg:
		return s.handleErrorOccuredMsg(msg)
	case events.LogEntriesLoadedMsg:
		return s.handleLogEntriesLoadedMsg(msg)
	case events.ViewRowsReloadRequestedMsg:
		return s.handleViewRowsReloadRequestedMsg()
	case events.OpenJSONRowRequestedMsg:
		return s.handleOpenJSONRowRequestedMsg(msg, s)
	case events.BackKeyClickedMsg:
		return s, tea.Quit
	case events.EnterKeyClickedMsg, events.ArrowRightKeyClickedMsg:
		return s.handleRequestOpenJSON()
	case events.FilterKeyClickedMsg:
		return s.handleFilterKeyClickedMsg()
	case tea.KeyMsg:
		cmdBatch = append(cmdBatch, s.handleKeyMsg(msg)...)

		if s.isFilterKeyMap(msg) {
			// Intercept table update.
			return s, tea.Batch(cmdBatch...)
		}
	}

	s.table, cmdBatch = batched(s.table.Update(msg))(cmdBatch)

	return s, tea.Batch(cmdBatch...)
}

func (s StateLoaded) handleKeyMsg(msg tea.KeyMsg) []tea.Cmd {
	var cmdBatch []tea.Cmd

	cmdBatch = appendCmd(cmdBatch, s.helper.handleKeyMsg(msg))

	if s.isArrowUpKeyMap(msg) {
		cmdBatch = appendCmd(cmdBatch, s.handleArrowUpKeyClicked())
	}

	return cmdBatch
}

func (s StateLoaded) handleArrowUpKeyClicked() tea.Cmd {
	if s.table.Cursor() == 0 {
		return events.ViewRowsReloadRequested
	}

	return nil
}

func (s StateLoaded) handleRequestOpenJSON() (tea.Model, tea.Cmd) {
	return s, events.OpenJSONRowRequested(s.logEntries, s.table.Cursor())
}

func (s StateLoaded) handleViewRowsReloadRequestedMsg() (tea.Model, tea.Cmd) {
	return s, s.helper.LoadEntries
}

func (s StateLoaded) handleFilterKeyClickedMsg() (tea.Model, tea.Cmd) {
	return initializeModel(newStateFiltering(s.helper.Application, s))
}

func (s StateLoaded) withApplication(application Application) (state, tea.Cmd) {
	s.helper.Application = application

	var cmd tea.Cmd
	s.table, cmd = s.table.Update(s.helper.Application.LastWindowSize)

	return s, cmd
}

// String implements fmt.Stringer.
func (s StateLoaded) String() string {
	return modelValue(s)
}

func (s StateLoaded) Application() Application {
	return s.helper.Application
}
