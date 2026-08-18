package widgets

import (
	"bytes"
	"strings"

	fx "github.com/antonmedv/fx/pkg/model"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hedhyw/json-log-viewer/internal/keymap"
)

const themeFX = "1"

// nolint: gochecknoinits // Dependency requirnment.
func init() {
	fx.SetCurrentThemeByID(themeFX)
}

// JSONViewModel is a widget that shows a prettified JSON tree. It wraps
// the fx model and adds a help screen with all supported hotkeys.
type JSONViewModel struct {
	jsonView tea.Model

	keyMap keymap.KeyMap

	windowSize tea.WindowSizeMsg
	showHelp   bool
	// inputFocused is true while the search or the dig input of the
	// JSON view is focused, so all keys belong to that input.
	inputFocused bool
}

// NewJSONViewModel creates a new JSON view widget if a content is the correct json,
// or plain text view otherwise.
func NewJSONViewModel(
	content []byte,
	lastWindowSize tea.WindowSizeMsg,
	keyMap keymap.KeyMap,
) (tea.Model, tea.Cmd) {
	fxModel, err := fx.New(fx.Config{
		FileName: "",
		Source:   bytes.NewReader(content),
	})
	if err != nil {
		return NewPlainLogModel(string(content), lastWindowSize, keyMap)
	}

	jsonView, cmd := fxModel.Update(lastWindowSize)

	return JSONViewModel{
		jsonView:   jsonView,
		keyMap:     keyMap,
		windowSize: lastWindowSize,
	}, cmd
}

// Init implements tea.Model interface.
func (m JSONViewModel) Init() tea.Cmd { return nil }

// View implements tea.Model interface.
func (m JSONViewModel) View() string {
	if m.showHelp {
		return m.viewHelp()
	}

	return m.jsonView.View()
}

// Update implements tea.Model interface.
func (m JSONViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.windowSize = msg
	}

	if msg, ok := msg.(tea.KeyMsg); ok {
		var handled bool

		if m, handled = m.handleKeyMsg(msg); handled {
			return m, nil
		}
	}

	var cmd tea.Cmd

	m.jsonView, cmd = m.jsonView.Update(msg)

	return m, cmd
}

// handleKeyMsg processes the keys of the help screen. The last returned
// value reports if the key is handled and should not be passed further.
func (m JSONViewModel) handleKeyMsg(msg tea.KeyMsg) (JSONViewModel, bool) {
	if m.showHelp {
		if key.Matches(msg, m.keyMap.Exit) {
			return m, false
		}

		// Any other key closes the help.
		m.showHelp = false

		return m, true
	}

	if !m.inputFocused && key.Matches(msg, m.keyMap.ToggleFullHelp) {
		m.showHelp = true

		return m, true
	}

	m.inputFocused = isInputFocused(m.inputFocused, msg)

	return m, false
}

// isInputFocused tracks if the search or the dig input of the JSON view
// is focused, because it is not exposed by the JSON view itself.
func isInputFocused(focused bool, msg tea.KeyMsg) bool {
	if focused {
		return msg.Type != tea.KeyEnter && msg.Type != tea.KeyEscape
	}

	return key.Matches(msg, jsonViewInputKeys)
}

// nolint: gochecknoglobals // Constant definition.
var jsonViewInputKeys = key.NewBinding(key.WithKeys("/", "."))

func (m JSONViewModel) viewHelp() string {
	fxKeys := fx.GetKeyMap()

	bindings := []fxBinding{
		{keys: m.keyMap.Back.Keys(), description: "back to the logs"},
		{keys: fxKeys.Up.Keys(), description: "up"},
		{keys: fxKeys.Down.Keys(), description: "down"},
		{keys: fxKeys.PageUp.Keys(), description: "page up"},
		{keys: fxKeys.PageDown.Keys(), description: "page down"},
		{keys: fxKeys.HalfPageUp.Keys(), description: "half page up"},
		{keys: fxKeys.HalfPageDown.Keys(), description: "half page down"},
		{keys: fxKeys.GotoTop.Keys(), description: "go to top"},
		{keys: fxKeys.GotoBottom.Keys(), description: "go to bottom"},
		{keys: fxKeys.NextSibling.Keys(), description: "next sibling"},
		{keys: fxKeys.PrevSibling.Keys(), description: "previous sibling"},
		{keys: fxKeys.Expand.Keys(), description: "expand"},
		{keys: fxKeys.Collapse.Keys(), description: "collapse"},
		{keys: fxKeys.ExpandRecursively.Keys(), description: "expand recursively"},
		{keys: fxKeys.CollapseRecursively.Keys(), description: "collapse recursively"},
		{keys: fxKeys.ExpandAll.Keys(), description: "expand all"},
		{keys: fxKeys.CollapseAll.Keys(), description: "collapse all"},
		{keys: fxKeys.ToggleWrap.Keys(), description: "toggle strings wrap"},
		{keys: fxKeys.Search.Keys(), description: "search regexp"},
		{keys: fxKeys.SearchNext.Keys(), description: "next search result"},
		{keys: fxKeys.SearchPrev.Keys(), description: "previous search result"},
		{keys: fxKeys.Preview.Keys(), description: "preview string"},
		{keys: fxKeys.Dig.Keys(), description: "dig json"},
		{keys: fxKeys.Yank.Keys(), description: "yank/copy (then y, k or p)"},
		{keys: m.keyMap.ToggleFullHelp.Keys(), description: "show/hide this help"},
		{keys: m.keyMap.Exit.Keys(), description: "exit"},
	}

	rendered := make([]string, 0, len(bindings))

	for _, b := range bindings {
		rendered = append(rendered, b.String())
	}

	return lipgloss.NewStyle().
		Width(m.windowSize.Width).
		MaxHeight(m.windowSize.Height).
		Render(strings.Join(rendered, "\n"))
}

// fxBinding is a single line of the help screen.
type fxBinding struct {
	keys        []string
	description string
}

// String implements fmt.Stringer.
func (b fxBinding) String() string {
	const keysWidth = 24

	keys := make([]string, 0, len(b.keys))

	for _, k := range b.keys {
		if k == " " {
			k = "space"
		}

		keys = append(keys, k)
	}

	return lipgloss.NewStyle().Width(keysWidth).Render(strings.Join(keys, ", ")) + b.description
}
