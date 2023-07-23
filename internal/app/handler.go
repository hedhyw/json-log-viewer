package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

const cellIDLogLevel = 1

func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		return m.back()
	case "q", "ctrl+c":
		return m.quit()
	case "enter":
		return m.handleEnter()
	case "f":
		return m.showFilter()
	default:
		return nil, nil
	}
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	if m.IsFilterShown() {
		return m.applyFilter()
	}

	if m.IsTableShown() || m.IsJSONShown() {
		return m.toggleLogEntity()
	}

	return nil, nil
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
	}

	m.table.SetRows(msg.Rows())
	m.filteredLogEntries = msg

	return m
}

func (m Model) handleErrorMsg(err error) Model {
	m.err = err

	return m
}
