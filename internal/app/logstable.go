package app

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

type logsTableModel struct {
	helper

	lazyTable      lazyTableModel[source.LazyLogEntry]
	lastWindowSize tea.WindowSizeMsg

	logEntries source.LazyLogEntries
}

func newLogsTableModel(application Application, logEntries source.LazyLogEntries) logsTableModel {
	helper := helper{Application: application}

	const cellIDLogLevel = 1

	tableLogs := table.New(
		table.WithColumns(getColumns(application.LastWindowSize.Width, application.Config)),
		table.WithFocused(true),
		table.WithHeight(application.LastWindowSize.Height),
	)

	tableLogs.SetStyles(getTableStyles())

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

	lazyTable := lazyTableModel[source.LazyLogEntry]{
		helper:          helper,
		table:           tableLogs,
		minRenderedRows: application.Config.PrerenderRows,
		allEntries:      logEntries,
		lastCursor:      0,
		renderedRows:    make([]table.Row, 0, application.Config.PrerenderRows*2),
	}.withRenderedRows()

	return logsTableModel{
		helper:     helper,
		lazyTable:  lazyTable,
		logEntries: logEntries,
	}.handleWindowSizeMsg(application.LastWindowSize)
}

// Init initializes component. It implements tea.Model.
func (m logsTableModel) Init() tea.Cmd {
	return m.lazyTable.Init()
}

// View renders component. It implements tea.Model.
func (m logsTableModel) View() string {
	return m.lazyTable.View()
}

// Update handles events. It implements tea.Model.
func (m logsTableModel) Update(msg tea.Msg) (logsTableModel, tea.Cmd) {
	var cmdBatch []tea.Cmd

	m.helper = m.helper.Update(msg)

	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m = m.handleWindowSizeMsg(msg)
	}

	m.lazyTable, cmdBatch = batched(m.lazyTable.Update(msg))(cmdBatch)

	return m, tea.Batch(cmdBatch...)
}

func (m logsTableModel) handleWindowSizeMsg(msg tea.WindowSizeMsg) logsTableModel {
	const heightOffset = 4

	x, y := m.BaseStyle.GetFrameSize()
	m.lazyTable.table.SetWidth(msg.Width - x*2)
	m.lazyTable.table.SetHeight(msg.Height - y*2 - footerSize - heightOffset)
	m.lazyTable.table.SetColumns(getColumns(m.lazyTable.table.Width()-10, m.Config))
	m.lastWindowSize = msg

	return m
}

// Cursor returns the index of the selected row.
func (m logsTableModel) Cursor() int {
	return m.lazyTable.table.Cursor()
}
