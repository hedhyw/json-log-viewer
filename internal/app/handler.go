package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

const cellIDLogLevel = 1

func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	for _, key := range m.table.KeyMap.LineUp.Keys() {
		if msg.String() == key {
			return m.handleUp()
		}
	}

	switch msg.String() {
	case "esc":
		return m.back()
	case "ctrl+c":
		return m.quit()
	case "q":
		return m.quitByRune()
	case "enter":
		return m.handleEnter()
	case "f":
		return m.showFilter()
	default:
		return nil, nil
	}
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	if m.IsErrorShown() {
		return m.quit()
	}

	if m.IsFilterShown() {
		return m.applyFilter()
	}

	if m.IsTableShown() || m.IsJSONShown() {
		return m.toggleLogEntity()
	}

	return nil, nil
}

func (m Model) handleUp() (tea.Model, tea.Cmd) {
	if m.table.Cursor() != 0 || !m.IsTableShown() || m.IsFiltered() {
		return nil, nil
	}

	m.allLogEntries = nil

	return m, m.Init()
}

func (m Model) handleWindowSizeMsg(msg tea.WindowSizeMsg) Model {
	x, y := m.baseStyle.GetFrameSize()
	m.table.SetWidth(msg.Width - x*2)
	m.table.SetHeight(msg.Height - y*2 - footerSize)
	m.table.SetColumns(getColumns(m.table.Width() - 10))
	m.lastWindowSize = msg

	return m
}

func (m Model) handleLogEntriesMsg(msg source.LogEntries) Model {
	if len(m.allLogEntries) == 0 {
		m.allLogEntries = msg
	}

	m.table.SetRows(msg.Rows())
	m.filteredLogEntries = msg

	tableStyles := getTableStyles()
	tableStyles.RenderCell = func(value string, rowID, columnID int) string {
		style := tableStyles.Cell

		if columnID == cellIDLogLevel {
			return removeClearSequence(
				m.getLogLevelStyle(style, rowID).Render(value),
			)
		}

		return style.Render(value)
	}

	m.table.SetStyles(tableStyles)

	m.table.UpdateViewport()

	return m
}

func (m Model) handleErrorMsg(err error) Model {
	m.err = err

	return m
}
