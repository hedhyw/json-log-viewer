package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hedhyw/bubbles/key"

	"github.com/hedhyw/json-log-viewer/internal/keymap"
	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
	"github.com/hedhyw/json-log-viewer/internal/pkg/widgets"
)

// StateFilteringModel is a state to prompt for filter term.
type StateFilteringModel struct {
	*Application

	previousState StateLoadedModel
	table         logsTableModel

	textInput widgets.PillInputModel
	keys      keymap.KeyMap

	// err holds a rejected filter term, it is shown next to the input so
	// that the user can correct the term without losing it.
	err error
}

func newStateFiltering(
	previousState StateLoadedModel,
) StateFilteringModel {
	fieldTitles := make([]string, len(previousState.Config.Fields))
	for i, f := range previousState.Config.Fields {
		fieldTitles[i] = f.Title
	}

	textInput := widgets.NewPillInputModel(fieldTitles)
	textInput.Focus()

	return StateFilteringModel{
		Application: previousState.Application,

		previousState: previousState,
		table:         previousState.table,

		textInput: textInput,
		keys:      previousState.getApplication().keys,
	}.resizeTable()
}

// footerSize is the number of lines rendered below the table: the input
// itself, its hints and an optional error.
func (s StateFilteringModel) footerSize() int {
	const inputSize = 2

	if s.err != nil {
		return inputSize + 1
	}

	return inputSize
}

// resizeTable fits the table into the space that is left by the footer.
func (s StateFilteringModel) resizeTable() StateFilteringModel {
	if size := s.footerSize(); s.table.footerSize != size {
		s.table.footerSize = size
		s.table = s.table.handleWindowSizeMsg(s.table.lastWindowSize)
	}

	return s
}

// Init initializes component. It implements tea.Model.
func (s StateFilteringModel) Init() tea.Cmd {
	return s.textInput.Init()
}

// View renders component. It implements tea.Model.
func (s StateFilteringModel) View() string {
	view := s.BaseStyle.Render(s.table.View()) + "\n" + s.textInput.View()

	if s.err != nil {
		view += "\n" + s.FooterStyle.Render(s.err.Error())
	}

	return view
}

// Update handles events. It implements tea.Model.
func (s StateFilteringModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmdBatch []tea.Cmd

	s.Application.Update(msg)

	switch msg := msg.(type) {
	case events.ErrorOccuredMsg:
		return s.handleErrorOccuredMsg(msg)
	case tea.KeyMsg:
		// Any key press means the user is amending the term, so a previously
		// reported problem with it is no longer relevant.
		s.err = nil
		s = s.resizeTable()

		if mdl, cmd := s.handleKeyMsg(msg); mdl != nil {
			return mdl, cmd
		}
	default:
		s.table, cmdBatch = batched(s.table.Update(msg))(cmdBatch)
	}

	var cmd tea.Cmd
	s.textInput, cmd = s.textInput.Update(msg)
	if cmd != nil {
		cmdBatch = append(cmdBatch, cmd)
	}

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
	filterField, input := s.textInput.Value()
	if input == "" {
		return s, events.EscKeyClicked
	}

	// Reject a malformed term here, while the user still has it in the input
	// and can correct it. Reporting it later would surface it as a fatal
	// application error.
	if _, err := source.NewMatcher(input); err != nil {
		s.err = err

		return s.resizeTable(), nil
	}

	return initializeModel(newStateFiltered(
		s.previousState,
		input,
		filterField,
	))
}

// String implements fmt.Stringer.
func (s StateFilteringModel) String() string {
	return modelValue(s)
}
