package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	defaultFooter = "[Ctrl+C] Exit; [Esc] Back; [Enter] Open/Hide; [↑↓] Navigation; [F] Filter"

	footerSize = 1
)

func (m Model) renderViews() string {
	if m.IsJSONShown() {
		return m.jsonView.View()
	}

	if m.IsErrorShown() {
		return fmt.Sprintf("something went wrong: %s", m.err)
	}

	var footer string

	if m.IsFilterShown() {
		footer = "\n" + m.textInput.View()
	} else {
		footer = "\n" + m.footerStyle.Render(defaultFooter)
	}

	return m.baseStyle.Render(m.table.View()) + footer
}

func (m Model) handleUpdateInViews(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	cmds := make([]tea.Cmd, 0, 3)

	if m.IsJSONShown() {
		_, cmd := m.jsonView.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.IsTableShown() {
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.IsFilterShown() {
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
