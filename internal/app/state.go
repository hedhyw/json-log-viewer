package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type stateModel interface {
	tea.Model
	fmt.Stringer

	getApplication() *Application
	refresh() (stateModel, tea.Cmd)
}
