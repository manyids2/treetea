package layout

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Close        key.Binding
	ViewTasks    key.Binding
	ViewContexts key.Binding
	ViewProjects key.Binding
	ViewTags     key.Binding
	ViewHistory  key.Binding
}

var keys = keyMap{
	Close: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "close pane"),
	),
	ViewTasks: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "Tasks"),
	),
	ViewContexts: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "Contexts"),
	),
	ViewProjects: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Projects"),
	),
	ViewTags: key.NewBinding(
		key.WithKeys("+"),
		key.WithHelp("+", "+tags"),
	),
	ViewHistory: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "histOry"),
	),
}
