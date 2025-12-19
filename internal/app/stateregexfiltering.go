package app

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/keymap"
	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
)

// StateRegexFilteringModel is a state to prompt for regex filter pattern.
type StateRegexFilteringModel struct {
	*Application

	previousState StateLoadedModel
	table         logsTableModel

	textInput textinput.Model
	keys      keymap.KeyMap
}

func newStateRegexFiltering(
	previousState StateLoadedModel,
) StateRegexFilteringModel {
	textInput := textinput.New()
	textInput.Focus()
	textInput.Placeholder = "Enter regex pattern..."

	return StateRegexFilteringModel{
		Application: previousState.Application,

		previousState: previousState,
		table:         previousState.table,

		textInput: textInput,
		keys:      previousState.getApplication().keys,
	}
}

func (s StateRegexFilteringModel) Init() tea.Cmd {
	return nil
}

func (s StateRegexFilteringModel) View() string {
	return s.BaseStyle.Render(s.table.View()) + "\n" +
		"Regex filter: " + s.textInput.View()
}

func (s StateRegexFilteringModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (s StateRegexFilteringModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, s.keys.Back) && string(msg.Runes) != "q":
		return s.previousState.refresh()
	case key.Matches(msg, s.keys.Open):
		return s.handleEnterKeyClickedMsg()
	default:
		return nil, nil
	}
}

func (s StateRegexFilteringModel) handleEnterKeyClickedMsg() (tea.Model, tea.Cmd) {
	if s.textInput.Value() == "" {
		return s, events.EscKeyClicked
	}

	return initializeModel(newStateRegexFiltered(
		s.previousState,
		s.textInput.Value(),
	))
}

func (s StateRegexFilteringModel) getApplication() *Application {
	return s.Application
}

func (s StateRegexFilteringModel) refresh() (stateModel, tea.Cmd) {
	return s, nil
}

func (s StateRegexFilteringModel) String() string {
	return modelValue(s)
}
