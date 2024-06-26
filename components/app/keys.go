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
	Project key.Binding
	Help    key.Binding

	Done        key.Binding
	StartStop   key.Binding
	Edit        key.Binding
	Editor      key.Binding
	Modify      key.Binding
	AddChild    key.Binding
	AddSibling  key.Binding
	Delete      key.Binding
	Select      key.Binding
	SelectTree  key.Binding
	SelectClear key.Binding

	ShowID      key.Binding
	ShowUUID    key.Binding
	ShowDue     key.Binding
	ShowTags    key.Binding
	ShowProject key.Binding
	ShowContext key.Binding

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
	Project: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "project"),
	),
	Accept: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "accept context"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel context"),
	),

	// Task actions
	Done: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "toggle done"),
	),
	StartStop: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "start/stop"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Editor: key.NewBinding(
		key.WithKeys("E"),
		key.WithHelp("E", "edit in editor"),
	),
	Modify: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "modify"),
	),
	AddChild: key.NewBinding(
		key.WithKeys("A"),
		key.WithHelp("A", "add child"),
	),
	AddSibling: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add sibling"),
	),
	Delete: key.NewBinding(
		key.WithKeys("ctrl+x"),
		key.WithHelp("ctrl+x", "delete task"),
	),
	Select: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "select"),
	),
	SelectTree: key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "select tree"),
	),
	SelectClear: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "clear selection"),
	),

	// Toggles
	ShowContext: key.NewBinding( // TODO: needs to be implemented
		key.WithKeys("C"),
		key.WithHelp("C", "show context"),
	),
	ShowTags: key.NewBinding(
		key.WithKeys("T"),
		key.WithHelp("T", "show tags"),
	),
	ShowDue: key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "show due"),
	),
	ShowID: key.NewBinding(
		key.WithKeys("I"),
		key.WithHelp("I", "show id"),
	),
	ShowUUID: key.NewBinding(
		key.WithKeys("U"),
		key.WithHelp("U", "show uuid"),
	),
	ShowProject: key.NewBinding(
		key.WithKeys("P"),
		key.WithHelp("P", "show project"),
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
	return []key.Binding{k.Quit, k.Help}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Help, k.Filter, k.Context, k.Project, k.ShowTags, k.ShowDue},
		{k.Up, k.Down, k.Top, k.Bottom, k.Left, k.Right},
		{k.Done, k.Edit, k.Editor, k.Modify, k.AddChild, k.AddSibling, k.Delete},
		{k.Select, k.SelectTree, k.SelectClear},
	}
}
