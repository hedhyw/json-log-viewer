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
	table := newLogsTableModel(application, logEntries)

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

func (s StateLoadedModel) viewHelp() string {
	toggles := func() string {
		toggles := []string{}
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

	if s.help.ShowAll {
		toggleText := lipgloss.NewStyle().
			Background(lipgloss.Color("#353533")).
			Padding(0, 1).
			Render(toggles())

		versionText := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#6124DF")).
			Padding(0, 1).
			Render(s.Version)

		width := s.Application.LastWindowSize.Width
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
	return "\n" + s.help.View(s.keys) + " " + toggles()
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
	return s, events.OpenJSONRowRequested(s.Entries, s.table.Cursor())
}

func (s StateLoadedModel) handleFilterKeyClickedMsg() (tea.Model, tea.Cmd) {
	return initializeModel(newStateFiltering(s))
}

func (s StateLoadedModel) getApplication() *Application {
	return s.Application
}

func (s StateLoadedModel) refresh() (stateModel, tea.Cmd) {
	var cmd tea.Cmd
	s.table, cmd = s.table.Update(s.Application.LastWindowSize)
	return s, cmd
}

// String implements fmt.Stringer.
func (s StateLoadedModel) String() string {
	return modelValue(s)
}
