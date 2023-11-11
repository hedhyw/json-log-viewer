package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// StateError is a failure message state.
type StateError struct {
	helper

	err error
}

func newStateError(application Application, err error) StateError {
	return StateError{
		helper: helper{Application: application},
		err:    err,
	}
}

// Init initializes component. It implements tea.Model.
func (s StateError) Init() tea.Cmd {
	return nil
}

// View renders component. It implements tea.Model.
func (s StateError) View() string {
	return fmt.Sprintf("Something went wrong: %s.", s.err)
}

// Update handles events. It implements tea.Model.
func (s StateError) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	s.helper = s.helper.Update(msg)

	switch msg.(type) {
	case tea.KeyMsg:
		return s, tea.Quit
	default:
		return s, nil
	}
}

// String implements fmt.Stringer.
func (s StateError) String() string {
	return modelValue(s)
}
