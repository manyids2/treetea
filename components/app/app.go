package app

import (
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/manyids2/tasktea/components/navbar"
	"github.com/manyids2/tasktea/components/tree"
	tk "github.com/manyids2/tasktea/task"
)

type editorFinishedMsg struct{ err error }

func openEditor(uuid string) tea.Cmd {
	c := exec.Command("task", "edit", uuid)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

type State int

const (
	StateHome = iota
	StateFilter
	StateContext
	StateProject
	StateEdit
	StateModify
	StateAdd
	StateInit
)

type Model struct {
	// Input
	Filters []string

	// Context
	Project      string
	Context      string
	ReadFilters  string
	WriteFilters string
	Active       string

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

	// Lists for filters
	contexts tree.Model
	projects tree.Model
	filters  tree.Model
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
		State:    StateInit,
		Width:    80,
		Height:   40,
		Padding:  "  ",
		Styles:   NewStyles(),
		nav:      navbar.New(),
		tree:     tree.New(),
		help:     help.New(),
		keys:     keys,
		contexts: tree.New(),
		projects: tree.New(),
		filters:  tree.New(),
	}

	// How to add default passdown args?
	m.tree.Padding = m.Padding
	m.nav.Padding = m.Padding

	// Get all available contexts
	m.RunContexts()
	m.RunProjects()

	// Get initial context, tasks
	m.RunContext()
	m.RunFilters()
	m.RunActive()

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

// RunContext Run `task context show` to get context, filters
func (m *Model) RunContext() {
	m.Context, m.ReadFilters, m.WriteFilters = tk.Context()
	m.SetProject()
	m.nav.Load(m.Project, m.Context, m.Filters)
}

// GetProject Find if `project:...` is part of filters
func (m *Model) SetProject() {
	project := ""
	parts := strings.Split(m.ReadFilters, " ")
	parts = append(parts, m.Filters...)
	for _, v := range parts {
		if len(v) < 9 {
			continue
		}
		if v[:8] == "project:" {
			project = v[8:]
		}
	}
	m.Project = project
}

// RunContext Run `task export active` to get context, filters
func (m *Model) RunActive() {
	m.Active = tk.Active()
	m.tree.Active = m.Active
}

func (m *Model) RunFilters() {
	read_filters := strings.Split(m.ReadFilters, " ")
	read_filters = append(read_filters, m.Filters...)
	read_filters = append(read_filters, "project:"+m.Project)
	tasks, _ := tk.List(read_filters)
	items := tasksToItems(tasks)
	m.tree.Load(items)
}

func (m *Model) RunContexts() {
	contexts := tk.Contexts()
	items := contextsToItems(contexts)
	m.contexts.Load(items)
}

func (m *Model) RunProjects() {
	projects := tk.Projects()
	items := contextsToItems(projects)
	m.projects.Load(items)
}

func (m Model) View() string {
	if m.quitting || m.State == StateInit {
		return ""
	}
	var content string
	switch m.State {
	case StateContext:
		content = m.nav.View() + "\n" +
			m.contexts.View() + "\n" +
			m.help.View(m.keys)
	case StateProject:
		content = m.nav.View() + "\n" +
			m.projects.View() + "\n" +
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

func (m Model) toggleStartStop(msg tea.Msg) (Model, tea.Cmd) {
	task := m.tree.CurrentItem().(tk.Task)
	if task.UUID == m.Active {
		tk.SetActive(task.UUID, "stop")
		m.Active = ""
		m.RunActive()
	} else {
		tk.SetActive(task.UUID, "start")
		m.RunActive()
	}
	return m, nil
}

func (m Model) deleteTask(msg tea.Msg) (Model, tea.Cmd) {
	task := m.tree.CurrentItem().(tk.Task)
	tk.Delete(task.UUID)
	m.RunFilters()
	return m, nil
}

func (m Model) editorTask(msg tea.Msg) (Model, tea.Cmd) {
	task := m.tree.CurrentItem().(tk.Task)
	cmd := openEditor(task.UUID)
	return m, cmd
}

func (m *Model) toggleExtra(extra string) {
	if m.tree.Extra == "" {
		m.tree.Extra = extra
	} else {
		m.tree.Extra = ""
	}
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

		// Project
		case key.Matches(msg, keys.Project):
			m.State = StateProject
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

		// Editor
		case key.Matches(msg, keys.Editor):
			return m.editorTask(msg)

		// Modify
		case key.Matches(msg, keys.Modify):
			m.State = StateModify

		// Add child
		case key.Matches(msg, keys.AddChild):
			m.State = StateAdd

		// Add sibling
		case key.Matches(msg, keys.AddSibling):
			m.State = StateAdd

		// Delete
		case key.Matches(msg, keys.Delete):
			return m.deleteTask(msg)

		// Toggle tags
		case key.Matches(msg, keys.ShowTags):
			m.toggleExtra("tags")
			return m, nil

		// Toggle due
		case key.Matches(msg, keys.ShowDue):
			m.toggleExtra("due")
			return m, nil

		// Toggle id
		case key.Matches(msg, keys.ShowID):
			m.toggleExtra("id")
			return m, nil

		// Toggle uuid
		case key.Matches(msg, keys.ShowUUID):
			m.toggleExtra("uuid")
			return m, nil

		// Toggle project
		case key.Matches(msg, keys.ShowProject):
			m.toggleExtra("project")
			return m, nil

		// Toggle uuid
		case key.Matches(msg, keys.StartStop):
			return m.toggleStartStop(msg)

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

func (m Model) selectProject(msg tea.Msg) (Model, tea.Cmd) {
	project := string(m.projects.CurrentItem().(tk.ContextS))

	// Load project as ReadFilters
	m.Project = project
	tk.SetContext("none")
	m.Context = "none"
	m.ReadFilters = ""
	m.WriteFilters = ""
	m.nav.Load(m.Project, m.Context, m.Filters)
	m.RunFilters()
	return m, nil
}

func (m Model) handleProject(msg tea.Msg) (Model, tea.Cmd) {
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
		case key.Matches(msg, keys.Project):
			m.State = StateHome
			return m, nil

		// Accept context
		case key.Matches(msg, m.keys.Accept):
			m.State = StateHome
			return m.selectProject(msg)

		// Dont allow edit
		case key.Matches(msg, m.keys.Edit):
			return m, nil

		// Cancel context
		case key.Matches(msg, m.keys.Cancel):
			m.State = StateHome
			return m, nil
		}
	}
	m.projects, cmd = m.projects.Update(msg) // Handle keys for contexts
	return m, cmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.nav.Width = msg.Width
		m.tree.Width = msg.Width
		m.tree.Height = msg.Height - 20 // Help is looong
		if m.State == StateInit {
			m.State = StateHome
		}

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

	case tree.ModifyCancelledMsg:
		m.State = StateHome

	case tree.ModifyChangedMsg:
		m.State = StateHome
		args := strings.Split(string(msg), " ")
		// If nothing is selected, modify current
		if len(m.tree.Selected) == 0 {
			task := m.tree.CurrentItem().(tk.Task)
			tk.Modify(task.UUID, args)
		} else {
			// Modify all selected
			for _, t := range m.tree.SelectedItems() {
				tk.Modify(t.(tk.Task).UUID, args)
			}
		}
		m.RunFilters() // Basically just reload

	case tree.AddCancelledMsg:
		m.State = StateHome

	case tree.AddChangedMsg:
		m.State = StateHome
		out := []string(msg)
		desc := out[0]
		parent := out[1]
		write_filters := strings.Split(m.WriteFilters, " ")
		write_filters = append(write_filters, m.Filters...)
		tk.Add(write_filters, desc, parent)
		m.RunFilters() // Basically just reload
	}

	switch m.State {
	case StateHome:
		m, cmd = m.handleHome(msg)
	case StateContext:
		m, cmd = m.handleContext(msg)
	case StateProject:
		m, cmd = m.handleProject(msg)
	case StateFilter:
		m.nav, cmd = m.nav.Update(msg) // Delegate to nav
	case StateEdit:
		m.tree, cmd = m.tree.Update(msg) // Delegate to tree
	case StateModify:
		m.tree, cmd = m.tree.Update(msg) // Delegate to tree
	case StateAdd:
		m.tree, cmd = m.tree.Update(msg) // Delegate to tree
	default: // should never occur
	}

	return m, cmd
}
