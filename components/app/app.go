package app

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	ly "github.com/manyids2/tasktea/components/layout"
	tw "github.com/manyids2/tasktea/task"
	xn "github.com/manyids2/tasktea/task/actions"
)

type (
	errMsg error
)

type Model struct {
	taskdata string
	taskrc   string
	rc       tw.TaskRC

	logpath  string
	quitting bool
	err      error

	Width  int
	Height int

	Context  string
	Filters  tw.Filters
	Active   []string
	Projects map[string]tw.Project
	Contexts map[string]tw.Filters

	Tasks []tw.Task

	layout ly.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func New() (m Model) {
	m = Model{
		layout: ly.New(),
	}
	m.LoadRc()
	m.LoadTasks()
	return m
}

func (m *Model) LoadRc() {
	var err error

	// Load taskrc
	m.taskdata, m.taskrc = xn.Paths()
	m.rc, err = tw.LoadTaskRC(m.taskrc)
	if err != nil {
		log.Fatalln("Could not load taskrc", err)
	}

	// Load available contexts with filters from rc
	m.Contexts = m.rc.Contexts

	// Get current context
	context, read_filters, write_filters, err := xn.Context()
	if err != nil {
		log.Fatalln("Could not retrieve context", err)
	}
	m.Context = context
	m.Filters = tw.Filters{
		Read:  read_filters,
		Write: write_filters,
	}

	// Get projects as tree
	m.Projects, err = xn.Projects()
	if err != nil {
		log.Fatalln("Could not retrieve projects", err)
	}

	// Check active
	m.Active, err = xn.Active()
	if err != nil {
		log.Fatalln("Could not retrieve active", err)
	}
}

func (m *Model) LoadTasks() {
	filters := strings.Split(m.Filters.Read, " ")
	tasks, err := xn.List(filters)
	if err != nil {
		log.Fatalln("Could not retrieve tasks", err)
	}
	m.Tasks = tasks
}
