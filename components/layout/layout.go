package layout

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
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
	ViewFilter
	ViewSave
)

const (
	TreeTasks = iota
	TreeContexts
	TreeProjects
	TreeTags
	TreeHistory
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

	// Filter
	filter textinput.Model
	save   textinput.Model

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

		filter: textinput.New(),
		save:   textinput.New(),
	}
	m.filter.Prompt = ""
	m.filter.Placeholder = ""
	m.filter.SetValue("")

	m.save.Prompt = ""
	m.save.Placeholder = ""
	m.save.SetValue("")
	return m
}

func (m Model) tree() *tr.Model {
	switch m.State {
	case ViewTasks:
		return &m.tasks
	case ViewProjects:
		return &m.projects
	case ViewContexts:
		return &m.contexts
	case ViewTags:
		return &m.tags
	case ViewHistory:
		return &m.history

	// TODO: There are really not in the same class as the others
	case ViewFilter:
		return &m.tasks
	case ViewSave:
		return &m.contexts
	}
	return &m.tasks
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

func (m Model) viewNav(tree *tr.Model) string {
	name := tree.Name
	desc := ""

	// Change only in case of tasks
	switch m.State {
	case ViewTasks:
		name = m.Context
		desc = m.Filters.Read
	case ViewFilter:
		// TODO: handle properly for each state
		name = m.Context
		desc = m.filter.View()
	case ViewSave:
		name = m.save.View()
		desc = m.Filters.Read
	}
	return fmt.Sprintf("| %s | %s\n", name, desc)
}

func (m Model) View() string {
	if !m.ready {
		return "not ready"
	}

	return m.frame.Render("\n" +
		m.viewNav(m.tree()) +
		"\n" +
		m.tree().View())
}

type errMsg error
type CloseMsg string
type ReloadRcMsg string

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
			m.State = ViewTasks
			m.LoadTasks(m.Context, m.Filters)
			// NOTE: Not changing tw context
			// xn.SetContext(context)
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	switch m.State {

	case ViewFilter:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.Cancel):
				m.filter.SetValue(m.Filters.Read)
				m.filter.Placeholder = m.Filters.Read
				m.State = ViewTasks
			case key.Matches(msg, m.keys.Accept):
				m.Filters.Read = m.filter.Value()
				m.State = ViewTasks
				m.LoadTasks(m.Context, m.Filters)
			}
		}
		m.filter, cmd = m.filter.Update(msg)

	case ViewSave:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.Cancel):
				m.State = ViewTasks
			case key.Matches(msg, m.keys.Accept):
				xn.SaveContext(m.save.Value(), m.Filters.Read)
				m.Context = m.save.Value()
				m.State = ViewTasks
				return m, func() tea.Msg { return ReloadRcMsg("") }
			}
		}
		m.save, cmd = m.save.Update(msg)

	default:
		// Key messages ( user input )
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			// Close
			case key.Matches(msg, m.keys.Close):
				return m, func() tea.Msg {
					return CloseMsg("")
				}
			// Change view
			case key.Matches(msg, m.keys.ViewTasks):
				m.State = ViewTasks
				return m, nil
			case key.Matches(msg, m.keys.ViewContexts):
				m.State = ViewContexts
				return m, nil
			case key.Matches(msg, m.keys.ViewHistory):
				m.State = ViewHistory
				return m, nil
			case key.Matches(msg, m.keys.ViewProjects):
				m.State = ViewProjects
				return m, nil
			case key.Matches(msg, m.keys.ViewTags):
				m.State = ViewTags
				return m, nil
			case key.Matches(msg, m.keys.Filter):
				m.State = ViewFilter
				m.filter.Placeholder = m.Filters.Read + " "
				m.filter.SetValue(m.Filters.Read + " ")
				m.filter.Focus()
				return m, nil
			case key.Matches(msg, m.keys.Save):
				m.State = ViewSave
				m.save.Placeholder = m.Context
				m.save.SetValue(m.Context)
				m.save.Focus()
				return m, nil

			// Toggle
			case key.Matches(msg, keys.ShowTags):
				m.tasks.ToggleExtra("tags")
				return m, nil
			case key.Matches(msg, keys.ShowDue):
				m.tasks.ToggleExtra("due")
				return m, nil
			case key.Matches(msg, keys.ShowID):
				m.tasks.ToggleExtra("id")
				return m, nil
			case key.Matches(msg, keys.ShowUUID):
				m.tasks.ToggleExtra("uuid")
				return m, nil
			case key.Matches(msg, keys.ShowProject):
				m.tasks.ToggleExtra("project")
				return m, nil
			}
		}

		// TODO: How to make this work with &tree?
		switch m.State {
		case ViewTasks:
			m.tasks, cmd = m.tasks.Update(msg)
		case ViewProjects:
			m.projects, cmd = m.projects.Update(msg)
		case ViewContexts:
			m.contexts, cmd = m.contexts.Update(msg)
		case ViewTags:
			m.tags, cmd = m.tags.Update(msg)
		case ViewHistory:
			m.history, cmd = m.history.Update(msg)

		// TODO: There are really not in the same class as the others
		case ViewFilter:
			m.tasks, cmd = m.tasks.Update(msg)
		case ViewSave:
			m.tasks, cmd = m.tasks.Update(msg)
		}
	}

	return m, cmd
}
