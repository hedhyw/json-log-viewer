package widgets

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hedhyw/bubbles/help"
	"github.com/hedhyw/bubbles/key"
	"github.com/hedhyw/bubbles/textinput"
)

// PillInputModel is a widget that allows inputting text with optional prefix selected from autocomplete suggestion.
type PillInputModel struct {
	textInput     textinput.Model
	help          help.Model
	keymap        inputKeymap
	filterField   string   // e.g., "level"
	isPillVisible bool     // true if a field is selected
	suggestions   []string // pill suggestions
}

// NewPillInputModel initializes a new PillInputModel with the given text.
// It updates a widget with the message `tea.WindowSizeMsg`.
func NewPillInputModel(suggestions []string) PillInputModel {
	ti := textinput.New()
	ti.Placeholder = "Field name or search term..."
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	ti.Focus()
	ti.ShowSuggestions = true
	ti.SetSuggestions(suggestions)

	h := help.New()

	km := inputKeymap{}

	return PillInputModel{textInput: ti, help: h, keymap: km, isPillVisible: false, filterField: "", suggestions: suggestions}
}

type inputKeymap struct{}

func (k inputKeymap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "complete")),
		key.NewBinding(key.WithKeys("down", "ctrl+n"), key.WithHelp("(↓, ctrl+n)", "next")),
		key.NewBinding(key.WithKeys("up", "ctrl+p"), key.WithHelp("(↑, ctrl+p)", "prev")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "quit")),
	}
}

func (k inputKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

func (m PillInputModel) Init() tea.Cmd {
	return nil
}

func (m PillInputModel) Update(msg tea.Msg) (PillInputModel, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyTab:
			if m.textInput.ShowSuggestions && len(m.textInput.AvailableSuggestions()) > 0 {
				m.filterField = m.textInput.CurrentSuggestion()
				m.isPillVisible = true

				// 1. Reset the input buffer completely
				m.textInput.Reset()
				m.textInput.SetValue("")
				m.textInput.SetSuggestions([]string{})
				m.textInput.Placeholder = "Search term..."

				// 2. Disable suggestions for the "Value" phase
				m.textInput.ShowSuggestions = false

				// 3. IMPORTANT: Return nil to stop the tab key from
				// being passed to the textinput.Update(msg) below.
				return m, nil
			}

		case tea.KeyBackspace:
			// If input is empty and a pill is active, jump back to field selection
			if m.textInput.Value() == "" && m.isPillVisible {
				m.isPillVisible = false
				m.textInput.Reset()
				m.textInput.SetValue(m.filterField)
				m.textInput.CursorEnd()
				m.filterField = "" // Clear the filter field when going back

				// Re-enable suggestions for the "Field" phase
				m.textInput.ShowSuggestions = true
				m.textInput.SetSuggestions(m.suggestions)
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	// Only update the input if we didn't handle a mode-switch above
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

var pillStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("229")).
	Background(lipgloss.Color("57")). // Purple background
	Padding(0, 1).
	MarginRight(1).
	Bold(true)

func (m PillInputModel) View() string {
	var s strings.Builder

	if m.isPillVisible {
		// Render the Pill
		pill := pillStyle.Render(m.filterField + ":")
		s.WriteString(pill)
	}

	// Render the input
	s.WriteString(m.textInput.View())

	return s.String() + "\n" + m.help.View(m.keymap)
}

func (m PillInputModel) Focus() tea.Cmd {
	return m.textInput.Focus()
}

func (m PillInputModel) Value() (string, string) {
	return m.filterField, m.textInput.Value()
}
