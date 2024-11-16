package keymap

import "github.com/charmbracelet/bubbles/key"

// KeyMap of the app.
type KeyMap struct {
	Exit            key.Binding
	Back            key.Binding
	Open            key.Binding
	ToggleViewArrow key.Binding
	Up              key.Binding
	Reverse         key.Binding
	Down            key.Binding
	Filter          key.Binding
	ToggleFullHelp  key.Binding
	GotoTop         key.Binding
	GotoBottom      key.Binding
	ShowPreview     key.Binding
}

// GetDefaultKeys returns default KeyMap.
func GetDefaultKeys() KeyMap {
	return KeyMap{
		Exit: key.NewBinding(
			key.WithKeys("ctrl+c", "f10"),
			key.WithHelp("Ctrl+C", "Exit"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "q"),
			key.WithHelp("esc", "Back"),
		),
		Open: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Open"),
		),
		ToggleViewArrow: key.NewBinding(
			key.WithKeys("right"),
		),
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "Up"),
		),
		Reverse: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "Reverse"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "Down"),
		),
		Filter: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "Filter"),
		),
		ToggleFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "Help"),
		),
		GotoTop: key.NewBinding(
			key.WithKeys("home"),
			key.WithHelp("home", "go to start"),
		),
		GotoBottom: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("end", "go to end"),
		),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Back, k.Open, k.Up, k.Down, k.ToggleFullHelp,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Back, k.Open},
		{k.Filter, k.Reverse},
		{k.GotoTop, k.GotoBottom},
		{k.ToggleFullHelp, k.Exit},
	}
}
