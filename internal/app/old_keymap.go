package app

import tea "github.com/charmbracelet/bubbletea"

func (a Application) isQuitKeyMap(
	msg tea.KeyMsg,
) bool {
	switch msg.String() {
	case "ctrl+c", "f10":
		return true
	default:
		return false
	}
}

