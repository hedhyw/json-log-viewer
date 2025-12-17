package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

// StateRegexFilteredModel is a state that shows regex-filtered records.
type StateRegexFilteredModel struct {
	*Application

	previousState StateLoadedModel
	table         logsTableModel
	logEntries    source.LazyLogEntries

	regexPattern string
}

func newStateRegexFiltered(
	previousState StateLoadedModel,
	regexPattern string,
) StateRegexFilteredModel {
	return StateRegexFilteredModel{
		Application: previousState.Application,

		previousState: previousState,

		regexPattern: regexPattern,
	}
}

func (s StateRegexFilteredModel) Init() tea.Cmd {
	return func() tea.Msg {
		return &s
	}
}

func (s StateRegexFilteredModel) View() string {
	footer := s.FooterStyle.Render(
		fmt.Sprintf("regex filtered %d by: /%s/", s.logEntries.Len(), s.regexPattern),
	)

	return s.BaseStyle.Render(s.table.View()) + "\n" + footer
}

func (s StateRegexFilteredModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmdBatch []tea.Cmd

	s.Application.Update(msg)

	if _, ok := msg.(*StateRegexFilteredModel); ok {
		s, msg = s.handleStateRegexFilteredModel()
	}

	if _, ok := msg.(events.LogEntriesUpdateMsg); ok {
		return s, nil
	}

	switch typedMsg := msg.(type) {
	case events.ErrorOccuredMsg:
		return s.handleErrorOccuredMsg(typedMsg)
	case events.OpenJSONRowRequestedMsg:
		return s.handleOpenJSONRowRequestedMsg(typedMsg, s)
	case tea.KeyMsg:
		if mdl, cmd := s.handleKeyMsg(typedMsg); mdl != nil {
			return mdl, cmd
		}
	}

	s.table, cmdBatch = batched(s.table.Update(msg))(cmdBatch)

	return s, tea.Batch(cmdBatch...)
}

func (s StateRegexFilteredModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, s.keys.Back):
		return s.previousState.refresh()
	case key.Matches(msg, s.keys.FilterRegex):
		return s.handleRegexFilterKeyClickedMsg()
	case key.Matches(msg, s.keys.ToggleViewArrow), key.Matches(msg, s.keys.Open):
		return s.handleRequestOpenJSON()
	default:
		return nil, nil
	}
}

func (s StateRegexFilteredModel) handleStateRegexFilteredModel() (StateRegexFilteredModel, tea.Msg) {
	entries, err := s.Application.Entries().FilterRegExp(s.regexPattern)
	if err != nil {
		return s, events.ShowError(err)()
	}

	s.logEntries = entries
	s.table = newLogsTableModel(
		s.Application,
		entries,
		false, // follow.
		s.previousState.table.lazyTable.reverse,
	)

	return s, nil
}

func (s StateRegexFilteredModel) handleRegexFilterKeyClickedMsg() (tea.Model, tea.Cmd) {
	state := newStateRegexFiltering(s.previousState)
	return initializeModel(state)
}

func (s StateRegexFilteredModel) handleRequestOpenJSON() (tea.Model, tea.Cmd) {
	if s.logEntries.Len() == 0 {
		return s, events.EscKeyClicked
	}

	return s, events.OpenJSONRowRequested(s.logEntries, s.table.Cursor())
}

func (s StateRegexFilteredModel) getApplication() *Application {
	return s.Application
}

func (s StateRegexFilteredModel) refresh() (_ stateModel, cmd tea.Cmd) {
	s.table, cmd = s.table.Update(s.LastWindowSize())

	return s, cmd
}

func (s StateRegexFilteredModel) String() string {
	return modelValue(s)
}
