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

	Quit    key.Binding
	QuitQ   key.Binding
	Filter  key.Binding
	Context key.Binding
	Help    key.Binding

	Done key.Binding
	Edit key.Binding

	Accept key.Binding
	Cancel key.Binding
}

var keys = keyMap{
	// Navigation
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

	// Filters
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),

	// Context
	Context: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "context"),
	),
	Accept: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "accept"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),

	// Task actions
	Done: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "toggle done"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),

	// Global
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q/esc", "quit"),
	),
	QuitQ: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Done, k.Filter, k.Help, k.Quit},
		{k.Up, k.Down, k.Top, k.Bottom, k.Left, k.Right},
		{k.Context, k.Accept, k.Cancel},
	}
}
