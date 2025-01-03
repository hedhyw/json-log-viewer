package app

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/keymap"
	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
)

// StateFilteringModel is a state to prompt for filter term.
type StateFilteringModel struct {
	*Application

	previousState StateLoadedModel
	table         logsTableModel

	textInput textinput.Model
	keys      keymap.KeyMap
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
		keys:      previousState.getApplication().keys,
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
		if mdl, cmd := s.handleKeyMsg(msg); mdl != nil {
			return mdl, cmd
		}
	default:
		s.table, cmdBatch = batched(s.table.Update(msg))(cmdBatch)
	}

	s.textInput, cmdBatch = batched(s.textInput.Update(msg))(cmdBatch)

	return s, tea.Batch(cmdBatch...)
}

func (s StateFilteringModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, s.keys.Back) && string(msg.Runes) != "q":
		return s.previousState.refresh()
	case key.Matches(msg, s.keys.Open):
		return s.handleEnterKeyClickedMsg()
	default:
		return nil, nil
	}
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
