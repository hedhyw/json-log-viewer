package app

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Exit       key.Binding
	Back       key.Binding
	BackQ       key.Binding
	ToggleView key.Binding
	ToggleViewArrow key.Binding
	Up         key.Binding
	Down       key.Binding
	Filter     key.Binding
}


var defaultKeys = KeyMap{
	Exit: key.NewBinding(
		key.WithKeys("ctrl", "c"),
		key.WithHelp("Ctrl+C", "Exit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("Esc", "Back"),
	),
	BackQ: key.NewBinding(
		key.WithKeys("q"),
	),
	ToggleView: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "Open/Hide"),
	),
	ToggleViewArrow: key.NewBinding(
		key.WithKeys("right"),
	),
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "Up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "Down"),
	),
	Filter: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("F", "Filter"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Exit, k.Back, k.ToggleView, k.Up, k.Down, k.Filter}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Up, k.Down},           // first column
		{k.ToggleView, k.Exit, k.Filter}, // second column
	}
}