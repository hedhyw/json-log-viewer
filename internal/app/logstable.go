package app

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

type logsTableModel struct {
	helper

	table          table.Model
	lastWindowSize tea.WindowSizeMsg

	logEntries source.LogEntries
}

func newLogsTableModel(application Application, logEntries source.LogEntries) logsTableModel {
	helper := helper{Application: application}

	const cellIDLogLevel = 1

	tableLogs := table.New(
		table.WithColumns(getColumns(application.LastWindowSize.Width, application.Config)),
		table.WithFocused(true),
		table.WithHeight(application.LastWindowSize.Height),
	)

	tableLogs.SetStyles(getTableStyles())
	tableLogs.SetRows(logEntries.Rows())

	tableStyles := getTableStyles()
	tableStyles.RenderCell = func(_ table.Model, value string, position table.CellPosition) string {
		style := tableStyles.Cell

		if position.Column == cellIDLogLevel {
			return removeClearSequence(
				helper.getLogLevelStyle(
					logEntries,
					style,
					position.RowID,
				).Render(value),
			)
		}

		return style.Render(value)
	}

	tableLogs.SetStyles(tableStyles)

	return logsTableModel{
		helper:     helper,
		table:      tableLogs,
		logEntries: logEntries,
	}.handleWindowSizeMsg(application.LastWindowSize)
}

// Init initializes component. It implements tea.Model.
func (m logsTableModel) Init() tea.Cmd {
	return nil
}

// View renders component. It implements tea.Model.
func (m logsTableModel) View() string {
	return m.table.View()
}

// Update handles events. It implements tea.Model.
func (m logsTableModel) Update(msg tea.Msg) (logsTableModel, tea.Cmd) {
	var cmdBatch []tea.Cmd

	m.helper = m.helper.Update(msg)

	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m = m.handleWindowSizeMsg(msg)
	}

	m.table, cmdBatch = batched(m.table.Update(msg))(cmdBatch)

	return m, tea.Batch(cmdBatch...)
}

func (m logsTableModel) handleWindowSizeMsg(msg tea.WindowSizeMsg) logsTableModel {
	const heightOffset = 4

	x, y := m.BaseStyle.GetFrameSize()
	m.table.SetWidth(msg.Width - x*2)
	m.table.SetHeight(msg.Height - y*2 - footerSize - heightOffset)
	m.table.SetColumns(getColumns(m.table.Width()-10, m.Config))
	m.lastWindowSize = msg

	return m
}

// Cursor returns the index of the selected row.
func (m logsTableModel) Cursor() int {
	return m.table.Cursor()
}
