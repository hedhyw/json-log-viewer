package widgets

import (
	"bytes"
	"encoding/json"

	fxjson "github.com/antonmedv/fx/pkg/json"
	fx "github.com/antonmedv/fx/pkg/model"
	"github.com/antonmedv/fx/pkg/theme"
	tea "github.com/charmbracelet/bubbletea"
)

const themeFX = "1"

// NewJSONViewModel creates a new JSON view widget if a content is the correct json,
// or plain text view otherwise.
func NewJSONViewModel(content []byte, lastWindowSize tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	jsonDecoder := json.NewDecoder(bytes.NewReader(content))
	jsonDecoder.UseNumber()

	object, err := fxjson.Parse(jsonDecoder)
	if err != nil {
		return NewPlainLogModel(string(content), lastWindowSize)
	}

	fxModel := fx.NewModel(object, fx.Config{
		FileName: "",
		Theme:    theme.Themes[themeFX],
		ShowSize: true,
	})

	return fxModel.Update(lastWindowSize)
}
