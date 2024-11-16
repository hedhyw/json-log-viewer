package widgets

import (
	"bytes"

	fx "github.com/antonmedv/fx/pkg/model"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/keymap"
)

const themeFX = "1"

// nolint: gochecknoinits // Dependency requirnment.
func init() {
	fx.SetCurrentThemeByID(themeFX)
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

	return fxModel.Update(lastWindowSize)
}
