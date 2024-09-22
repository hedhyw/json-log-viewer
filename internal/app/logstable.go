package app

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/events"

	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

type logsTableModel struct {
	*Application

	lazyTable      lazyTableModel
	lastWindowSize tea.WindowSizeMsg
	footerSize     int

	logEntries source.LazyLogEntries
}

func newLogsTableModel(application *Application, logEntries source.LazyLogEntries) logsTableModel {
	const cellIDLogLevel = 1

	tableLogs := table.New(
		table.WithColumns(getColumns(application.LastWindowSize.Width, application.Config)),
		table.WithFocused(true),
		table.WithHeight(application.LastWindowSize.Height),
	)
	tableLogs.KeyMap.LineUp = application.keys.Up
	tableLogs.KeyMap.LineDown = application.keys.Down
	tableLogs.KeyMap.GotoBottom = application.keys.GotoBottom
	tableLogs.KeyMap.GotoTop = application.keys.GotoTop

	tableLogs.SetStyles(getTableStyles())

	tableStyles := getTableStyles()
	tableStyles.RenderCell = func(_ table.Model, value string, position table.CellPosition) string {
		style := tableStyles.Cell

		if position.Column == cellIDLogLevel {
			return removeClearSequence(
				application.getLogLevelStyle(
					logEntries,
					style,
					position.RowID,
				).Render(value),
			)
		}

		return style.Render(value)
	}

	tableLogs.SetStyles(tableStyles)

	lazyTable := lazyTableModel{
		Application:  application,
		reverse:      true,
		follow:       true,
		table:        tableLogs,
		entries:      logEntries,
		lastCursor:   0,
		renderedRows: nil,
	}

	msg := logsTableModel{
		Application: application,
		lazyTable:   lazyTable,
		logEntries:  logEntries,
		footerSize:  1,
	}.handleWindowSizeMsg(application.LastWindowSize)

	return msg
}

// View renders component. It implements tea.Model.
func (m logsTableModel) View() string {
	return m.lazyTable.View()
}

// Update handles events. It implements tea.Model.
func (m logsTableModel) Update(msg tea.Msg) (logsTableModel, tea.Cmd) {
	var cmdBatch []tea.Cmd

	m.Application.Update(msg)

	switch typedMsg := msg.(type) {
	case tea.WindowSizeMsg:
		m = m.handleWindowSizeMsg(typedMsg)
	case events.LogEntriesUpdateMsg:
		m.logEntries = source.LazyLogEntries(typedMsg)
		msg = EntriesUpdateMsg{Entries: m.logEntries}
	}

	m.lazyTable, cmdBatch = batched(m.lazyTable.Update(msg))(cmdBatch)

	return m, tea.Batch(cmdBatch...)
}

func (m logsTableModel) handleWindowSizeMsg(msg tea.WindowSizeMsg) logsTableModel {
	const (
		heightOffset = 4
		widthOffset  = -10
	)

	x, y := m.BaseStyle.GetFrameSize()
	m.lazyTable.table.SetWidth(msg.Width - x*2)
	m.lazyTable.table.SetHeight(msg.Height - y*2 - m.footerSize - heightOffset)
	m.lazyTable.table.SetColumns(getColumns(m.lazyTable.table.Width()+widthOffset, m.Config))
	m.lastWindowSize = msg

	m.lazyTable = m.lazyTable.RenderedRows()

	return m
}

// Cursor returns the index of the selected row.
func (m logsTableModel) Cursor() int {
	return m.lazyTable.viewPortCursor()
}
