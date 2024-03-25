package layout

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"

	nv "github.com/manyids2/tasktea/components/navbar"
	st "github.com/manyids2/tasktea/components/statusbar"
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
	frame  lg.Style
	keys   keyMap

	// State of layout
	State   ViewState
	Context string
	Filters tw.Filters
	Tasks   []tw.Task

	// Information
	nav    nv.Model
	status st.Model

	// Trees
	tasks    tr.Model
	contexts tr.Model
	projects tr.Model
	tags     tr.Model
	history  tr.Model

	// Error handling
	ready bool
	err   error
}

func New() (m Model) {
	m = Model{
		Width:  80,
		Height: 24,
		State:  ViewTasks,
		nav:    nv.New(),
		status: st.New(),
		keys:   keys,
	}
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

	// Update UI
	m.State = ViewTasks
	m.nav.Title = m.Context
	m.nav.Description = m.Filters.Read
	m.status.Message = fmt.Sprintf("%d Tasks", len(m.Tasks))
}

func (m *Model) Layout() {
	m.nav.Width = m.Width
	m.nav.Height = 3

	m.status.Width = m.Width
	m.status.Height = 2
}

func (m *Model) SetFrame(width, height int) {
	m.Width, m.Height = width, height
	m.frame = lg.NewStyle().Height(height).Width(width)
}

func (m Model) View() string {
	if !m.ready {
		return "not ready"
	}
	switch m.State {
	case ViewTasks:
		m.nav.Title = m.Context
	case ViewContexts:
		m.nav.Title = "Contexts"
	case ViewProjects:
		m.nav.Title = "Projects"
	case ViewTags:
		m.nav.Title = "Tags"
	case ViewHistory:
		m.nav.Title = "History"
	}
	return m.frame.Render(m.nav.View() + "\n" + m.status.View())
}

type errMsg error

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	// Non-key messages ( event listeners )
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		// m.Height = msg.Height // Dont change height for now
		m.Layout()
		m.ready = true
		return m, nil

	case errMsg:
		m.err = msg
		return m, nil
	}

	// Key messages ( user input )
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.ViewTasks):
			m.State = ViewTasks
			return m, nil
		case key.Matches(msg, m.keys.ViewProjects):
			m.State = ViewProjects
			return m, nil
		case key.Matches(msg, m.keys.ViewContexts):
			m.State = ViewContexts
			return m, nil
		case key.Matches(msg, m.keys.ViewTags):
			m.State = ViewTags
			return m, nil
		case key.Matches(msg, m.keys.ViewHistory):
			m.State = ViewHistory
			return m, nil
		}
	}

	return m, cmd
}
