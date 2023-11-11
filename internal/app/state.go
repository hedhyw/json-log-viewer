package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type state interface {
	tea.Model
	fmt.Stringer

	withApplication(application Application) (state, tea.Cmd)
}
