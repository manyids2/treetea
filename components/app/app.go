package app

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/manyids2/tasktea/components/navbar"
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
	Filters []string

	// Context
	Context      string
	ReadFilters  string
	WriteFilters string

	// State
	State    State
	quitting bool

	// Appearance
	Width   int
	Height  int
	Padding string
	Styles  Styles

	// Components
	nav  navbar.Model
	tree tree.Model
	help help.Model
	keys keyMap
}

func tasksToItems(tasks []tk.Task) []tree.Item {
	items := []tree.Item{}
	for _, t := range tasks {
		items = append(items, t)
	}
	return items
}

func NewModel(filters []string) Model {
	m := Model{
		Filters: filters,
		State:   StateHome,
		Width:   80,
		Height:  1,
		Padding: "  ",
		Styles:  NewStyles(),
		nav:     navbar.New(),
		tree:    tree.New(),
		help:    help.New(),
		keys:    keys,
	}

	// How to add default passdown args?
	m.tree.Padding = m.Padding
	m.nav.Padding = m.Padding

	// Get initial context, tasks
	m.Context, m.ReadFilters, m.WriteFilters = tk.Context()
	m.nav.Load(m.Context, m.Filters)
	m.RunFilters()

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) RunFilters() {
	read_filters := strings.Split(m.ReadFilters, " ")
	read_filters = append(read_filters, m.Filters...)
	tasks, _ := tk.List(read_filters)
	items := tasksToItems(tasks)
	m.tree.Load(items)
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	return m.nav.View() + "\n" +
		m.tree.View() + "\n" +
		m.help.View(m.keys)
}

func (m Model) toggleDone(msg tea.Msg) (Model, tea.Cmd) {
	task := m.tree.CurrentItem().(tk.Task)
	switch task.Status {
	case "pending":
		tk.SetStatus(task.UUID, "completed")
		m.RunFilters()
	case "completed":
		tk.SetStatus(task.UUID, "pending")
		m.RunFilters()
	}
	return m, nil
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

		case key.Matches(msg, m.keys.ToggleDone):
			m, cmd = m.toggleDone(msg)
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
		m.nav.Width = msg.Width
		m.tree.Width = msg.Width
		m.tree.Height = msg.Height - 10

	case navbar.CancelledMsg:
		m.State = StateHome

	case navbar.ChangedMsg:
		m.Filters = []string(msg)
		m.State = StateHome
		m.RunFilters()
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
