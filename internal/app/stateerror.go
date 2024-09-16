package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// StateErrorModel is a failure message state.
type StateErrorModel struct {
	*Application

	err error
}

func newStateError(application *Application, err error) StateErrorModel {
	return StateErrorModel{
		Application: application,
		err:         err,
	}
}

// Init initializes component. It implements tea.Model.
func (s StateErrorModel) Init() tea.Cmd {
	return nil
}

// View renders component. It implements tea.Model.
func (s StateErrorModel) View() string {
	return fmt.Sprintf("Something went wrong: %s.", s.err)
}

// Update handles events. It implements tea.Model.
func (s StateErrorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	s.Application.Update(msg)

	switch msg.(type) {
	case tea.KeyMsg:
		return s, tea.Quit
	default:
		return s, nil
	}
}

// String implements fmt.Stringer.
func (s StateErrorModel) String() string {
	return modelValue(s)
}
