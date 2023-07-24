package app

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/widgets"
)

func (m Model) hideJSON() (Model, tea.Cmd) {
	m.jsonView = nil

	return m, nil
}

func (m Model) showJSON() (Model, tea.Cmd) {
	cursor := m.table.Cursor()
	if cursor < 0 || cursor >= len(m.filteredLogEntries) {
		return m, nil
	}

	logEntry := m.filteredLogEntries[cursor]

	jsonViewModel, cmd := widgets.NewJSONViewModel(logEntry.Line, m.lastWindowSize)
	m.jsonView = jsonViewModel

	return m, cmd
}

func (m Model) quit() (Model, tea.Cmd) {
	return m, tea.Quit
}

func (m Model) quitByRune() (tea.Model, tea.Cmd) {
	if m.IsFilterShown() {
		return nil, nil
	}

	return m, tea.Quit
}

func (m Model) showFilter() (tea.Model, tea.Cmd) {
	if m.IsFilterShown() {
		return nil, nil
	}

	if !m.IsTableShown() {
		return nil, nil
	}

	m.textInputShown = true
	m.textInput = textinput.New()
	m.textInput.Focus()
	m.textInput.Prompt = "Filter >"
	m.table.Blur()

	return m, nil
}

func (m Model) toggleLogEntity() (Model, tea.Cmd) {
	if m.IsJSONShown() {
		return m.hideJSON()
	}

	return m.showJSON()
}

func (m Model) clearFilter() (Model, tea.Cmd) {
	m.textInput.SetValue("")

	return m.applyFilter()
}

func (m Model) back() (Model, tea.Cmd) {
	if m.IsFilterShown() {
		return m.clearFilter()
	}

	if m.IsJSONShown() {
		return m.hideJSON()
	}

	if m.IsFiltered() {
		return m.clearFilter()
	}

	return m.quit()
}

func (m Model) applyFilter() (Model, tea.Cmd) {
	m.textInputShown = false
	m.table.Focus()

	term := m.textInput.Value()

	if term == "" {
		return m, func() tea.Msg {
			return m.allLogEntries
		}
	}

	m.table.GotoTop()

	return m, func() tea.Msg {
		return m.allLogEntries.Filter(term)
	}
}
