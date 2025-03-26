package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

// StateLoadedModel is a state that shows all loaded records.
type StateLoadedModel struct {
	*Application

	table logsTableModel
}

func newStateViewLogs(
	application *Application,
	logEntries source.LazyLogEntries,
) StateLoadedModel {
	table := newLogsTableModel(
		application,
		logEntries,
		true, // follow.
		true, // reverse.
	)

	return StateLoadedModel{
		Application: application,

		table: table,
	}
}

// Init initializes component. It implements tea.Model.
func (s StateLoadedModel) Init() tea.Cmd {
	return nil
}

// View renders component. It implements tea.Model.
func (s StateLoadedModel) View() string {
	return s.viewTable() + s.viewHelp()
}

func (s StateLoadedModel) viewTable() string {
	return s.BaseStyle.Render(s.table.View())
}

func (s StateLoadedModel) toggles() string {
	toggles := make([]string, 0, 2)

	if s.table.lazyTable.reverse {
		toggles = append(toggles, "reverse")
	}

	if s.table.lazyTable.follow {
		toggles = append(toggles, "following")
	}

	if len(toggles) > 0 {
		return fmt.Sprintf("( %s )", strings.Join(toggles, ", "))
	}

	return ""
}

func (s StateLoadedModel) viewHelp() string {
	if s.help.ShowAll {
		toggleText := lipgloss.NewStyle().
			Background(lipgloss.Color("#353533")).
			Padding(0, 1).
			Render(s.toggles())

		versionText := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#6124DF")).
			Padding(0, 1).
			Render(s.Version)

		width := s.LastWindowSize().Width
		fillerText := lipgloss.NewStyle().
			Background(lipgloss.Color("#353533")).
			Width(width - lipgloss.Width(toggleText) - lipgloss.Width(versionText)).
			Render("")

		bar := lipgloss.JoinHorizontal(lipgloss.Top,
			toggleText,
			fillerText,
			versionText,
		)

		return "\n" + s.help.View(s.keys) + "\n" + lipgloss.NewStyle().Width(width).Render(bar)
	}
	return "\n" + s.help.View(s.keys) + " " + s.toggles()
}

// Update handles events. It implements tea.Model.
//
// nolint: cyclop // Many events in switch case.
func (s StateLoadedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmdBatch []tea.Cmd

	s.Application.Update(msg)

	switch msg := msg.(type) {
	case events.ErrorOccuredMsg:
		return s.handleErrorOccuredMsg(msg)
	case events.OpenJSONRowRequestedMsg:
		return s.handleOpenJSONRowRequestedMsg(msg, s)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.Back):
			return s, tea.Quit
		case key.Matches(msg, s.keys.Filter):
			return s.handleFilterKeyClickedMsg()
		case key.Matches(msg, s.keys.ToggleViewArrow), key.Matches(msg, s.keys.Open):
			return s.handleRequestOpenJSON()
		case key.Matches(msg, s.keys.ToggleFullHelp):
			s.help.ShowAll = !s.help.ShowAll
			if s.help.ShowAll {
				s.table.footerSize = 3
			} else {
				s.table.footerSize = 1
			}
			return s.refresh()
		}
		cmdBatch = append(cmdBatch, s.handleKeyMsg(msg)...)
	}

	s.table, cmdBatch = batched(s.table.Update(msg))(cmdBatch)

	return s, tea.Batch(cmdBatch...)
}

func (s StateLoadedModel) handleKeyMsg(msg tea.KeyMsg) []tea.Cmd {
	var cmdBatch []tea.Cmd

	cmdBatch = appendCmd(cmdBatch, s.Application.handleKeyMsg(msg))

	return cmdBatch
}

func (s StateLoadedModel) handleRequestOpenJSON() (tea.Model, tea.Cmd) {
	return s, events.OpenJSONRowRequested(s.entries, s.table.Cursor())
}

func (s StateLoadedModel) handleFilterKeyClickedMsg() (tea.Model, tea.Cmd) {
	return initializeModel(newStateFiltering(s))
}

func (s StateLoadedModel) getApplication() *Application {
	return s.Application
}

func (s StateLoadedModel) refresh() (_ stateModel, cmd tea.Cmd) {
	var cmdFirst, cmdSecond tea.Cmd

	s.table, cmdSecond = s.table.Update(s.LastWindowSize())
	s.table, cmdFirst = s.table.Update(events.LogEntriesUpdateMsg(s.Entries()))

	return s, tea.Batch(cmdFirst, cmdSecond)
}

// String implements fmt.Stringer.
func (s StateLoadedModel) String() string {
	return modelValue(s)
}
