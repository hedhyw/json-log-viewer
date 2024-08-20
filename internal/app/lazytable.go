package app

import (
	"slices"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
)

// rowGetter renders the row.
type rowGetter interface {
	// Row return a rendered table row.
	Row(cfg *config.Config, i int) table.Row
	Len() int
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		render = m.handleKey(msg, render)

	case EntriesUpdateMsg:
		m.entries = msg.Entries
		render = true
	}
	m.table, cmd = m.table.Update(msg)

	if m.table.Cursor() != m.lastCursor {
		render = true
	}

	if render {
		m = m.RenderedRows()
	}

	return m, cmd
}

func (m lazyTableModel) handleKey(msg tea.KeyMsg, render bool) bool {
	if key.Matches(msg, m.Application.keys.Reverse) {
		m.reverse = !m.reverse
		render = true
	}

	increaseOffset := func() {
		maxO := max(m.entries.Len()-m.table.Height(), 0)
		o := min(m.offset+1, maxO)
		if o != m.offset {
			m.offset = o
			render = true
		} else {
			m.follow = true
		}
	}
	decreaseOffset := func() {
		o := max(m.offset-1, 0)
		if o != m.offset {
			m.offset = o
			render = true
		}
	}
	if m.reverse {
		increaseOffset, decreaseOffset = decreaseOffset, increaseOffset
	}
	if key.Matches(msg, m.Application.keys.Down) {
		m.follow = false
		if m.table.Cursor()+1 == m.table.Height() {
			increaseOffset()
		}
	}
	if key.Matches(msg, m.Application.keys.Up) {
		m.follow = false
		if m.table.Cursor() == 0 {
			decreaseOffset()
		}
	}
	if key.Matches(msg, m.Application.keys.GotoTop) {
		if m.reverse {
			m.follow = true
		} else {
			m.follow = false
			m.offset = 0
		}
		render = true
	}
	if key.Matches(msg, m.Application.keys.GotoBottom) {
		if m.reverse {
			m.follow = false
			m.offset = 0
		} else {
			m.follow = true
		}
		render = true
	}

	return render
}

func (m lazyTableModel) ViewPortCursor() int {
	if m.reverse {
		viewSize := m.ViewPortEnd() - m.ViewPortStart()

		return m.offset + (viewSize - 1 - m.table.Cursor())
	}

	return m.offset + m.table.Cursor()
}

func (m lazyTableModel) ViewPortStart() int {
	return m.offset
}

func (m lazyTableModel) ViewPortEnd() int {
	return min(m.offset+m.table.Height(), m.entries.Len())
}

func (m lazyTableModel) RenderedRows() lazyTableModel {
	if m.follow {
		m.offset = max(0, m.entries.Len()-m.table.Height())
	}
	end := min(m.offset+m.table.Height(), m.entries.Len())

	m.renderedRows = []table.Row{}
	for i := m.offset; i < end; i++ {
		m.renderedRows = append(m.renderedRows, m.entries.Row(m.Config, i))
	}

	if m.reverse {
		slices.Reverse(m.renderedRows)
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

	return m
}
