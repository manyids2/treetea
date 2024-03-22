package app

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Top    key.Binding
	Bottom key.Binding

	Left  key.Binding
	Right key.Binding

	Quit   key.Binding
	Filter key.Binding
	Help   key.Binding

	ToggleDone key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Top: key.NewBinding(
		key.WithKeys("t", "g"),
		key.WithHelp("t/g", "move to top"),
	),
	Bottom: key.NewBinding(
		key.WithKeys("b", "G"),
		key.WithHelp("b/G", "move to bottom"),
	),
	Left: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("←/h", "prev page"),
	),
	Right: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("→/l", "next page"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q/esc", "quit"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	ToggleDone: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "toggle done"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Filter, k.Help, k.Quit},
		{k.Up, k.Down, k.Top, k.Bottom},
		{k.Left, k.Right},
		{k.ToggleDone},
	}
}
