package app

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

// StateLoadedModel is a state that shows all loaded records.
type StateLoadedModel struct {
	helper

	initCmd tea.Cmd

	table        logsTableModel
	logEntries   source.LazyLogEntries
	lastReloadAt time.Time

	keys      KeyMap
	help      help.Model
	reloading bool
}

func newStateViewLogs(
	application Application,
	logEntries source.LazyLogEntries,
	lastReloadAt time.Time,
) StateLoadedModel {
	table := newLogsTableModel(application, logEntries)

	return StateLoadedModel{
		helper: helper{Application: application},

		initCmd: table.Init(),

		table:      table,
		logEntries: logEntries,

		keys: defaultKeys,
		help: help.New(),

		lastReloadAt: lastReloadAt,
	}
}

// Init initializes component. It implements tea.Model.
func (s StateLoadedModel) Init() tea.Cmd {
	return s.initCmd
}

// View renders component. It implements tea.Model.
func (s StateLoadedModel) View() string {
	if s.reloading {
		return s.viewTable() + "\nreloading..."
	}

	return s.viewTable() + s.viewHelp()
}

func (s StateLoadedModel) viewTable() string {
	return s.BaseStyle.Render(s.table.View())
}

func (s StateLoadedModel) viewHelp() string {
	return "\n" + s.help.View(s.keys)
}

// Update handles events. It implements tea.Model.
//
// nolint: cyclop // Many events in switch case.
func (s StateLoadedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmdBatch []tea.Cmd

	s.helper = s.helper.Update(msg)

	switch msg := msg.(type) {
	case events.ErrorOccuredMsg:
		return s.handleErrorOccuredMsg(msg)
	case events.LogEntriesLoadedMsg:
		return s.handleLogEntriesLoadedMsg(msg, s.lastReloadAt)
	case events.ViewRowsReloadRequestedMsg:
		return s.handleViewRowsReloadRequestedMsg()
	case events.OpenJSONRowRequestedMsg:
		return s.handleOpenJSONRowRequestedMsg(msg, s)
	case tea.KeyMsg:
		if s.reloading {
			return s, nil
		}

		switch {
		case key.Matches(msg, s.keys.Back):
			return s, tea.Quit
		case key.Matches(msg, s.keys.Filter):
			return s.handleFilterKeyClickedMsg()
		case key.Matches(msg, s.keys.ToggleViewArrow), key.Matches(msg, s.keys.ToggleView):
			return s.handleRequestOpenJSON()
		}
		cmdBatch = append(cmdBatch, s.handleKeyMsg(msg)...)
	}

	s.table, cmdBatch = batched(s.table.Update(msg))(cmdBatch)

	return s, tea.Batch(cmdBatch...)
}

func (s StateLoadedModel) handleKeyMsg(msg tea.KeyMsg) []tea.Cmd {
	var cmdBatch []tea.Cmd

	cmdBatch = appendCmd(cmdBatch, s.helper.handleKeyMsg(msg))

	if key.Matches(msg, s.keys.Up) {
		cmdBatch = appendCmd(cmdBatch, s.handleArrowUpKeyClicked())
	}

	return cmdBatch
}

func (s StateLoadedModel) handleArrowUpKeyClicked() tea.Cmd {
	if s.table.Cursor() == 0 {
		return events.ViewRowsReloadRequested
	}

	return nil
}

func (s StateLoadedModel) handleRequestOpenJSON() (tea.Model, tea.Cmd) {
	if len(s.logEntries) == 0 {
		return s, tea.Quit
	}

	return s, events.OpenJSONRowRequested(s.logEntries, s.table.Cursor())
}

func (s StateLoadedModel) handleViewRowsReloadRequestedMsg() (tea.Model, tea.Cmd) {
	if time.Since(s.lastReloadAt) < s.Config.ReloadThreshold || s.reloading {
		return s, nil
	}

	s.lastReloadAt = time.Now()
	s.reloading = true

	return s, s.helper.LoadEntries
}

func (s StateLoadedModel) handleFilterKeyClickedMsg() (tea.Model, tea.Cmd) {
	return initializeModel(newStateFiltering(s.helper.Application, s))
}

func (s StateLoadedModel) withApplication(application Application) (stateModel, tea.Cmd) {
	s.helper.Application = application

	var cmd tea.Cmd
	s.table, cmd = s.table.Update(s.helper.Application.LastWindowSize)

	return s, cmd
}

// String implements fmt.Stringer.
func (s StateLoadedModel) String() string {
	return modelValue(s)
}

func (s StateLoadedModel) Application() Application {
	return s.helper.Application
}
