package app

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

// Model of the application.
type Model struct {
	baseStyle   lipgloss.Style
	footerStyle lipgloss.Style

	fileLogPath string

	table         table.Model
	allLogEntries source.LogEntries

	filteredLogEntries source.LogEntries

	lastWindowSize tea.WindowSizeMsg
	jsonView       tea.Model

	textInputShown bool
	textInput      textinput.Model

	err error
}

// NewModel initializes a new application model. It accept the path
// to the file with logs.
func NewModel(path string) Model {
	tableLogs := table.New(
		table.WithColumns(getColumns(100)),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	tableLogs.SetStyles(getTableStyles())

	return Model{
		baseStyle:   getBaseStyle(),
		footerStyle: getFooterStyle(),

		fileLogPath: path,
		table:       tableLogs,

		err:                nil,
		allLogEntries:      nil,
		filteredLogEntries: nil,

		textInputShown: false,
		textInput:      textinput.Model{},

		lastWindowSize: tea.WindowSizeMsg{},
		jsonView:       nil,
	}
}

// Init implements team.Model interface.
func (m Model) Init() tea.Cmd {
	return source.LoadLogsFromFile(m.fileLogPath)
}

// Update implements team.Model interface.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m = m.handleWindowSizeMsg(msg)
	case source.LogEntries:
		m = m.handleLogEntriesMsg(msg)
	case error:
		m = m.handleErrorMsg(msg)

		return m, nil
	case tea.KeyMsg:
		newModel, cmd := m.handleKeyMsg(msg)
		if newModel != nil || cmd != nil {
			return newModel, cmd
		}
	}

	return m.handleUpdateInViews(msg)
}

// View implements team.Model interface.
func (m Model) View() string {
	return m.renderViews()
}
