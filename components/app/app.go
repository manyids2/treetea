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
	StateContext
	StateEdit
	StateAdd
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

	// Main components
	nav  navbar.Model // Context, Filters
	tree tree.Model   // Tasks
	help help.Model
	keys keyMap

	contexts tree.Model // Context selection
}

func tasksToItems(tasks []tk.Task) []tree.Item {
	items := []tree.Item{}
	for _, t := range tasks {
		items = append(items, t)
	}
	return items
}

func contextsToItems(contexts []string) []tree.Item {
	items := []tree.Item{}
	for _, t := range contexts {
		items = append(items, tk.ContextS(t))
	}
	return items
}

func NewModel(filters []string) Model {
	m := Model{
		Filters:  filters,
		State:    StateHome,
		Width:    80,
		Height:   1,
		Padding:  "  ",
		Styles:   NewStyles(),
		nav:      navbar.New(),
		tree:     tree.New(),
		help:     help.New(),
		keys:     keys,
		contexts: tree.New(),
	}

	// How to add default passdown args?
	m.tree.Padding = m.Padding
	m.nav.Padding = m.Padding

	// Get all available contexts
	m.RunContexts()

	// Get initial context, tasks
	m.RunContext()
	m.RunFilters()

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) RunContext() {
	m.Context, m.ReadFilters, m.WriteFilters = tk.Context()
	m.nav.Load(m.Context, m.Filters)
}

func (m *Model) RunFilters() {
	read_filters := strings.Split(m.ReadFilters, " ")
	read_filters = append(read_filters, m.Filters...)
	tasks, _ := tk.List(read_filters)
	items := tasksToItems(tasks)
	m.tree.Load(items)
}

func (m *Model) RunContexts() {
	contexts := tk.Contexts()
	items := contextsToItems(contexts)
	m.contexts.Load(items)
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	var content string
	switch m.State {
	case StateContext:
		content = m.nav.View() + "\n" +
			m.contexts.View() + "\n" +
			m.help.View(m.keys)
	default:
		content = m.nav.View() + "\n" +
			m.tree.View() + "\n" +
			m.help.View(m.keys)
	}
	return content
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
		case key.Matches(msg, keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		// Context
		case key.Matches(msg, keys.Context):
			m.State = StateContext
			return m, nil

		// Filter
		case key.Matches(msg, keys.Filter):
			m.nav.StartFilter()
			m.State = StateFilter

		// Done
		case key.Matches(msg, keys.Done):
			return m.toggleDone(msg)

		// Edit
		case key.Matches(msg, keys.Edit):
			m.State = StateEdit

		// Add child
		case key.Matches(msg, keys.AddChild):
			m.State = StateAdd

		// Add sibling
		case key.Matches(msg, keys.AddSibling):
			m.State = StateAdd
		}
	}
	m.tree, cmd = m.tree.Update(msg) // Handle keys for tree
	return m, cmd
}

func (m Model) selectContext(msg tea.Msg) (Model, tea.Cmd) {
	context := string(m.contexts.CurrentItem().(tk.ContextS))
	tk.SetContext(context)
	m.RunContext()
	m.RunFilters()
	return m, nil
}

func (m Model) handleContext(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		// Quit only q, ctrl-c
		case key.Matches(msg, keys.QuitQ):
			m.quitting = true
			return m, tea.Quit

		// Help
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		// Context
		case key.Matches(msg, keys.Context):
			m.State = StateHome
			return m, nil

		// Accept context
		case key.Matches(msg, m.keys.Accept):
			m.State = StateHome
			return m.selectContext(msg)

		// Dont allow edit
		case key.Matches(msg, m.keys.Edit):
			return m, nil

		// Cancel context
		case key.Matches(msg, m.keys.Cancel):
			m.State = StateHome
			return m, nil
		}
	}
	m.contexts, cmd = m.contexts.Update(msg) // Handle keys for contexts
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

	case tree.EditCancelledMsg:
		m.State = StateHome

	case tree.EditChangedMsg:
		m.State = StateHome
		task := m.tree.CurrentItem().(tk.Task)
		desc := string(msg)
		tk.ModifyDescription(task.UUID, desc)
		m.RunFilters() // Basically just reload

	case tree.AddCancelledMsg:
		m.State = StateHome

	case tree.AddChangedMsg:
		m.State = StateHome
	}

	switch m.State {
	case StateHome:
		m, cmd = m.handleHome(msg)
	case StateContext:
		m, cmd = m.handleContext(msg)
	case StateFilter:
		m.nav, cmd = m.nav.Update(msg) // Delegate to nav
	case StateEdit:
		m.tree, cmd = m.tree.Update(msg) // Delegate to tree
	case StateAdd:
		m.tree, cmd = m.tree.Update(msg) // Delegate to tree
	default: // should never occur
	}

	return m, cmd
}
