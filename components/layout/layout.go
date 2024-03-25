package layout

import (
	"fmt"
	"log"
	"strings"

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
	ViewHome = iota
)

type Model struct {
	// Essential
	Width  int
	Height int
	State  ViewState

	// State of layout
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
		nav:    nv.New(),
		status: st.New(),
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

func (m Model) View() string {
	var style = lg.NewStyle().
		Height(m.Height).
		Width(m.Width)
	if !m.ready {
		return "not ready"
	}
	return style.Render(m.nav.View() + "\n" + m.status.View())
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

	return m, cmd
}
