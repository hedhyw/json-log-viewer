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

func (a Application) isEnterKeyMap(msg tea.KeyMsg) bool {
	return msg.String() == "enter"
}

func (a Application) isArrowUpKeyMap(msg tea.KeyMsg) bool {
	return msg.Type == tea.KeyUp
}

func (a Application) isArrowRightKeyMap(msg tea.KeyMsg) bool {
	return msg.Type == tea.KeyRight
}
