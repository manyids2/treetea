package app

import (
	"log"

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

	// From taskrc
	Contexts map[string]tw.Filters

	// Querying task
	Active   []string
	Tags     []string
	Projects map[string]tw.Project

	// Initial context and filters
	Context string
	Filters tw.Filters

	Width  int
	Height int

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
	m.layout.LoadTasks(m.Context, m.Filters)
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

	// Get tags as list
	m.Tags, err = xn.Tags()
	if err != nil {
		log.Fatalln("Could not retrieve projects", err)
	}

	// Check active
	m.Active, err = xn.Active()
	if err != nil {
		log.Fatalln("Could not retrieve active", err)
	}
}
