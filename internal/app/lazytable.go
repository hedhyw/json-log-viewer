package app

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
)

// rowGetter renders the row.
type rowGetter interface {
	// Row return a rendered table row.
	Row(cfg *config.Config) table.Row
}

// lazyTableModel lazily renders table rows.
type lazyTableModel[T rowGetter] struct {
	helper

	table table.Model

	minRenderedRows int
	allEntries      []T
	lastCursor      int

	renderedRows []table.Row
}

// Init implements tea.Model.
func (m lazyTableModel[T]) Init() tea.Cmd {
	return nil
}

// View implements tea.Model.
func (m lazyTableModel[T]) View() string {
	return m.table.View()
}

// Update implements tea.Model.
func (m lazyTableModel[T]) Update(msg tea.Msg) (lazyTableModel[T], tea.Cmd) {
	var cmd tea.Cmd

	m.table, cmd = m.table.Update(msg)

	if m.table.Cursor() != m.lastCursor {
		m = m.withRenderedRows()
	}

	return m, cmd
}

func (m lazyTableModel[T]) withRenderedRows() lazyTableModel[T] {
	cursor := m.table.Cursor()

	start := max(len(m.renderedRows), cursor)
	end := min(cursor+m.minRenderedRows, len(m.allEntries))

	for i := start; i < end; i++ {
		m.renderedRows = append(m.renderedRows, m.allEntries[i].Row(m.Config))
	}

	m.table.SetRows(m.renderedRows)
	m.lastCursor = m.table.Cursor()

	return m
}
