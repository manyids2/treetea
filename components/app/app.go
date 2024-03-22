package app

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/manyids2/tasktea/components/navbar"
	"github.com/manyids2/tasktea/components/statusbar"
	"github.com/manyids2/tasktea/components/tree"
	tk "github.com/manyids2/tasktea/task"
)

type State int

const (
	StateHome = iota
	StateFilter
)

type Model struct {
	// Input
	Context string
	Filters []string

	// State
	State    State
	quitting bool

	// Appearance
	Width   int
	Height  int
	Padding string
	Styles  Styles

	// Components
	nav    navbar.Model
	status statusbar.Model
	tree   tree.Model
	help   help.Model
	keys   keyMap
}

func NewModel(context string, filters []string) Model {
	// Get some tasks
	tasks, _ := tk.List(filters)
	items := []tree.Item{}
	for _, t := range tasks {
		items = append(items, t)
	}

	// items := []tree.Item{
	// 	tk.Task{UUID: "A"},
	// 	tk.Task{UUID: "B"},
	// 	tk.Task{UUID: "C", Depends: []string{"A"}},
	// }

	m := Model{
		Context: context,
		Filters: filters,
		State:   StateHome,
		Width:   80,
		Height:  1,
		Padding: "  ",
		Styles:  NewStyles(),
		nav:     navbar.New(context, filters),
		status:  statusbar.New(),
		tree:    tree.New(items),
		help:    help.New(),
		keys:    keys,
	}

	m.tree.Padding = m.Padding
	m.nav.Padding = m.Padding

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	return m.nav.View() + "\n" +
		m.tree.View() + "\n" +
		m.help.View(m.keys)
}

func (m Model) handleHome(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		// Quit
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			return m, tea.Quit

		// Help
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		// Edit filter
		case key.Matches(msg, keys.Filter):
			m.nav.StartFilter()
			m.State = StateFilter
		}
	}
	m.tree, cmd = m.tree.Update(msg) // Handle keys for tree
	return m, cmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width

	case navbar.CancelledMsg:
		m.State = StateHome

	case navbar.ChangedMsg:
		m.Filters = []string(msg)
		m.State = StateHome
	}

	switch m.State {
	case StateHome:
		m, cmd = m.handleHome(msg)
	case StateFilter:
		m.nav, cmd = m.nav.Update(msg) // Delegate to nav
	default: // should never occur
	}

	return m, cmd
}
