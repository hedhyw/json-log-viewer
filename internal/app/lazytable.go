package app

import (
	"slices"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

// rowGetter renders the row.
type rowGetter interface {
	// Row return a rendered table row.
	Row(cfg *config.Config, i int) table.Row
	// Len returns the number of all rows.
	Len() int
	// LogEntry getter
	LogEntry(cfg *config.Config, i int) source.LogEntry
}

// lazyTableModel lazily renders table rows.
type lazyTableModel struct {
	*Application

	table table.Model

	entries    rowGetter
	lastCursor int
	offset     int
	reverse    bool
	follow     bool

	renderedRows []table.Row
}

type EntriesUpdateMsg struct {
	Entries rowGetter
}

// View implements tea.Model.
func (m lazyTableModel) View() string {
	return m.table.View()
}

// Update implements tea.Model.
func (m lazyTableModel) Update(msg tea.Msg) (lazyTableModel, tea.Cmd) {
	var cmd tea.Cmd

	render := false
	captureMessage := false
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m, render, captureMessage = m.handleKey(msg, render)

	case EntriesUpdateMsg:
		m.entries = msg.Entries
		render = true
	}

	if !captureMessage {
		m.table, cmd = m.table.Update(msg)
	}

	if m.table.Cursor() != m.lastCursor {
		render = true
	}

	if render {
		m = m.RenderedRows()
	}

	return m, cmd
}

func (m lazyTableModel) getCellRenderer() func(table.Model, string, table.CellPosition) string {
	cellIDLogLevel := getIndexByKind(m.Config, config.FieldKindLevel)
	tableStyles := getTableStyles()

	return func(_ table.Model, value string, position table.CellPosition) string {
		style := tableStyles.Cell

		if position.Column == cellIDLogLevel {
			return removeClearSequence(
				m.getLogLevelStyle(
					m.renderedRows,
					style,
					position.RowID,
				).Render(value),
			)
		}

		return style.Render(value)
	}
}

func (m lazyTableModel) handleKey(msg tea.KeyMsg, render bool) (lazyTableModel, bool, bool) {
	captureMessage := false // when true, the key message must not be forwarded to the inner table

	// toggle the reverse display of items.
	if key.Matches(msg, m.keys.Reverse) {
		m.reverse = !m.reverse
		render = true
	}

	// this function increases the viewport offset by n if possible.  (scrolls down)
	increaseOffset := func(n int) {
		maxOffset := max(m.entries.Len()-m.table.Height(), 0)

		offset := min(m.offset+n, maxOffset)
		if offset != m.offset {
			m.offset = offset
			render = true
		} else {
			// we were at the last item, so we should follow the log
			m.follow = true
		}
	}

	// this function decreases the viewport offset by n if possible.  (scrolls up)
	decreaseOffset := func(n int) {
		offset := max(m.offset-n, 0)
		if offset != m.offset {
			m.offset = offset
			render = true
		}
	}

	// if the table is being displayed in reverse order, we need to swap the increase and decrease functions
	// since the last item is at the top of the table instead of the bottom.
	if m.reverse {
		increaseOffset, decreaseOffset = decreaseOffset, increaseOffset
	}

	if key.Matches(msg, m.keys.Down) {
		m.follow = false
		if m.table.Cursor()+1 == m.table.Height() {
			increaseOffset(1) // move the viewport
		}
	}

	if key.Matches(msg, m.keys.Up) {
		m.follow = false
		if m.table.Cursor() == 0 {
			decreaseOffset(1) // move the viewport
		}
	}

	if key.Matches(msg, m.keys.PageDown) {
		m.follow = false
		increaseOffset(m.table.Height() - 1) // move the viewport
		captureMessage = !m.follow
	}

	if key.Matches(msg, m.keys.PageUp) {
		m.follow = false
		decreaseOffset(m.table.Height() - 1) // move the viewport
		captureMessage = !m.follow
	}

	if key.Matches(msg, m.keys.GotoTop) {
		if m.reverse {
			// when follow is enabled, rendering will handle setting the offset to the correct value
			m.follow = true
		} else {
			m.follow = false
			m.offset = 0
		}
		render = true
	}

	if key.Matches(msg, m.keys.GotoBottom) {
		if m.reverse {
			m.follow = false
			m.offset = 0
		} else {
			// when follow is enabled, rendering will handle setting the offset to the correct value
			m.follow = true
		}
		render = true
	}

	return m, render, captureMessage
}

func (m lazyTableModel) viewPortCursor() int {
	if m.reverse {
		viewSize := m.viewPortEnd() - m.viewPortStart()

		return m.offset + (viewSize - 1 - m.table.Cursor())
	}

	return m.offset + m.table.Cursor()
}

func (m lazyTableModel) viewPortStart() int {
	return m.offset
}

func (m lazyTableModel) viewPortEnd() int {
	return min(m.offset+m.table.Height(), m.entries.Len())
}

// RenderedRows returns current visible rendered rows.
func (m lazyTableModel) RenderedRows() lazyTableModel {
	if m.follow {
		m.offset = max(0, m.entries.Len()-m.table.Height())
	}
	end := min(m.offset+m.table.Height(), m.entries.Len())

	m.renderedRows = m.renderedRows[:0]
	renderedEntries := make([]source.LogEntry, 0, cap(m.renderedRows))
	for i := m.offset; i < end; i++ {
		m.renderedRows = append(m.renderedRows, m.entries.Row(m.Config, i))
		renderedEntries = append(renderedEntries, m.entries.LogEntry(m.Config, i))
	}

	if m.reverse {
		slices.Reverse(m.renderedRows)
		slices.Reverse(renderedEntries)
	}

	m.table.SetRows(m.renderedRows)
	if m.follow {
		if m.reverse {
			m.table.GotoTop()
		} else {
			m.table.GotoBottom()
		}
	}

	m.lastCursor = m.table.Cursor()

	tableStyles := getTableStyles()
	tableStyles.RenderCell = m.getCellRenderer()
	m.table.SetStyles(tableStyles)

	return m
}
