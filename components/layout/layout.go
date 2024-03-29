package layout

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"

	tr "github.com/manyids2/tasktea/components/tree"
	tw "github.com/manyids2/tasktea/task"
	xn "github.com/manyids2/tasktea/task/actions"
)

type ViewState int

const (
	ViewTasks = iota
	ViewContexts
	ViewProjects
	ViewTags
	ViewHistory
)

type Model struct {
	// Essential
	Width  int
	Height int

	// State of layout
	State   ViewState
	Context string
	Filters tw.Filters
	Tasks   []tw.Task

	// Trees
	tasks    tr.Model
	contexts tr.Model
	projects tr.Model
	tags     tr.Model
	history  tr.Model

	// Current tree
	tree tr.Model

	// Helpers
	frame lg.Style
	keys  keyMap

	// Error handling
	ready bool
	err   error
}

func New() (m Model) {
	m = Model{
		Width:  80,
		Height: 24,
		State:  ViewTasks,
		keys:   keys,

		tasks:    tr.New("Tasks", "Tasks", "Tasks"),
		contexts: tr.New("Contexts", "Context", "Contexts"),
		projects: tr.New("Projects", "Project", "Projects"),
		tags:     tr.New("Tags", "Tag", "Tags"),
		history:  tr.New("History", "Item", "Items"),
	}
	m.tree = m.tasks
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) LoadTasks(context string, filters tw.Filters) {
	// Set context
	m.Context = context
	m.Filters = filters

	// Retrieve tasks with `task [filters] export`
	var args []string
	if m.Filters.Read != "" {
		args = strings.Split(m.Filters.Read, " ")
	} else {
		args = []string{}
	}
	tasks, err := xn.List(args)
	if err != nil {
		log.Fatalln("Could not retrieve tasks", err)
	}
	m.Tasks = tasks

	// Put tasks into items: NOTE: How to not copy???
	items := tr.Items{}
	for _, v := range m.Tasks {
		items = append(items, v)
	}
	m.tasks.LoadTree(items)

	// Update UI
	m.State = ViewTasks
	m.tree = m.tasks
}

func (m *Model) LoadContexts(contexts map[string]tw.Filters) {
	items := tr.Items{}
	for k, v := range contexts {
		items = append(items, tw.ContextItem{Name: k, Filters: v})
	}
	m.contexts.LoadList(items)
}

func (m *Model) LoadProjects(projects []string) {
	items := tr.Items{}
	for _, v := range projects {
		items = append(items, tw.StringItem(v))
	}
	m.projects.LoadList(items)
}

func (m *Model) LoadTags(tags []string) {
	items := tr.Items{}
	for _, v := range tags {
		items = append(items, tw.StringItem(v))
	}
	m.tags.LoadList(items)
}

func (m *Model) LoadHistory(history []string) {
	items := tr.Items{}
	for _, v := range history {
		items = append(items, tw.StringItem(v))
	}
	m.history.LoadList(items)
}

func (m *Model) SetFrame(width, height int) {
	m.Width, m.Height = width, height
	m.frame = lg.NewStyle().Height(height).Width(width)
}

func (m Model) viewNav() string {
	name := m.tree.Name
	desc := ""

	// Change only in case of tasks
	if m.State == ViewTasks {
		name = m.Context
		desc = m.Filters.Read
	}
	return fmt.Sprintf("| %s | %s\n", name, desc)
}

func (m Model) View() string {
	if !m.ready {
		return "not ready"
	}

	return m.frame.Render("\n" +
		m.viewNav() +
		"\n" +
		m.tree.View())
}

type errMsg error

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	// Non-key messages ( event listeners )
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.ready = true
		return m, nil

	case tr.AcceptMsg:
		switch m.State {
		case ViewContexts:
			context := msg[0]
			m.Context = context
			_, filters := m.contexts.Items.Get(context)
			m.Filters = tw.Filters{Read: filters.Desc("read"), Write: filters.Desc("write")}
			m.LoadTasks(m.Context, m.Filters)
			// NOTE: Not changing tw context
			// xn.SetContext(context)
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	// Key messages ( user input )
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		// Change view
		case key.Matches(msg, m.keys.ViewTasks):
			m.State = ViewTasks
			m.tree = m.tasks
			return m, nil
		case key.Matches(msg, m.keys.ViewContexts):
			m.State = ViewContexts
			m.tree = m.contexts
			return m, nil
		case key.Matches(msg, m.keys.ViewHistory):
			m.State = ViewHistory
			m.tree = m.history
			return m, nil
		case key.Matches(msg, m.keys.ViewProjects):
			m.State = ViewProjects
			m.tree = m.projects
			return m, nil
		case key.Matches(msg, m.keys.ViewTags):
			m.State = ViewTags
			m.tree = m.tags
			return m, nil
		}
	}

	m.tree, cmd = m.tree.Update(msg) // Delegate to current tree
	return m, cmd
}
