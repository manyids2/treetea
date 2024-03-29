package layout

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Close key.Binding

	ViewTasks    key.Binding
	ViewContexts key.Binding
	ViewProjects key.Binding
	ViewTags     key.Binding
	ViewHistory  key.Binding

	Filter key.Binding
	Accept key.Binding
	Cancel key.Binding
	Save   key.Binding
}

var keys = keyMap{
	Close: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "close pane"),
	),

	ViewTasks: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "tasks"),
	),
	ViewContexts: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "contexts"),
	),
	ViewProjects: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "projects"),
	),
	ViewTags: key.NewBinding(
		key.WithKeys("+"),
		key.WithHelp("+", "tags"),
	),
	ViewHistory: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "old"),
	),

	Filter: key.NewBinding(
		key.WithKeys("/", "f"),
		key.WithHelp("/", "filter"),
	),
	Accept: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "accept"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	Save: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "save"),
	),
}
