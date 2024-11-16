package widgets

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"

	"github.com/hedhyw/json-log-viewer/internal/keymap"
)

// PlainLogModel is a widget that shows multiline text in a viewport.
type PlainLogModel struct {
	viewport viewport.Model
	text     string
	keyMap   keymap.KeyMap
}

// NewPlainLogModel initializes a new PlainLogModel with the given text.
// It updates a widget with the message `tea.WindowSizeMsg`.
func NewPlainLogModel(
	text string,
	windowSize tea.WindowSizeMsg,
	keyMap keymap.KeyMap,
) (PlainLogModel, tea.Cmd) {
	m := PlainLogModel{
		text:     text,
		viewport: viewport.New(windowSize.Width, windowSize.Height),
		keyMap:   keyMap,
	}

	m = m.refreshText(windowSize.Width)

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(windowSize)

	return m, cmd
}

// Init implements team.Model interface.
func (m PlainLogModel) Init() tea.Cmd { return nil }

// View implements team.Model interface.
func (m PlainLogModel) View() string {
	return m.viewport.View()
}

// Update implements team.Model interface.
func (m PlainLogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// nolint: gocritic // For future extension.
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
		m = m.refreshText(msg.Width)
	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.Exit) ||
			key.Matches(msg, m.keyMap.Back) ||
			key.Matches(msg, m.keyMap.Open) {
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)

	return m, cmd
}

func (m PlainLogModel) refreshText(width int) PlainLogModel {
	m.viewport.SetContent(wordwrap.String(m.text, width))

	return m
}
