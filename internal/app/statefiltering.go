package app

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
)

// StateFilteringModel is a state to prompt for filter term.
type StateFilteringModel struct {
	*Application

	previousState StateLoadedModel
	table         logsTableModel

	textInput textinput.Model
	keys      KeyMap
}

func newStateFiltering(
	previousState StateLoadedModel,
) StateFilteringModel {
	textInput := textinput.New()
	textInput.Focus()

	return StateFilteringModel{
		Application: previousState.Application,

		previousState: previousState,
		table:         previousState.table,

		textInput: textInput,
		keys:      defaultKeys,
	}
}

// Init initializes component. It implements tea.Model.
func (s StateFilteringModel) Init() tea.Cmd {
	return nil
}

// View renders component. It implements tea.Model.
func (s StateFilteringModel) View() string {
	return s.BaseStyle.Render(s.table.View()) + "\n" + s.textInput.View()
}

// Update handles events. It implements tea.Model.
func (s StateFilteringModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmdBatch []tea.Cmd

	s.Application.Update(msg)

	switch msg := msg.(type) {
	case events.ErrorOccuredMsg:
		return s.handleErrorOccuredMsg(msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.Back):
			return s.previousState.refresh()
		case key.Matches(msg, s.keys.Open):
			return s.handleEnterKeyClickedMsg()
		}
		if cmd := s.handleKeyMsg(msg); cmd != nil {
			// Intercept table update.
			return s, cmd
		}
	default:
		s.table, cmdBatch = batched(s.table.Update(msg))(cmdBatch)
	}

	s.textInput, cmdBatch = batched(s.textInput.Update(msg))(cmdBatch)

	return s, tea.Batch(cmdBatch...)
}

func (s StateFilteringModel) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	if len(msg.Runes) == 1 {
		return nil
	}

	return s.Application.handleKeyMsg(msg)
}

func (s StateFilteringModel) handleEnterKeyClickedMsg() (tea.Model, tea.Cmd) {
	if s.textInput.Value() == "" {
		return s, events.EscKeyClicked
	}

	return initializeModel(newStateFiltered(
		s.previousState,
		s.textInput.Value(),
	))
}

// String implements fmt.Stringer.
func (s StateFilteringModel) String() string {
	return modelValue(s)
}
