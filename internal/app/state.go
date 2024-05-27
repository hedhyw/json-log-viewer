package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type stateModel interface {
	tea.Model
	fmt.Stringer

	withApplication(application Application) (stateModel, tea.Cmd)
}
