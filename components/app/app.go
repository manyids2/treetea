package app

import (
	"log"

	tw "github.com/manyids2/tasktea/task"
	"github.com/manyids2/tasktea/task/actions"
)

type Model struct {
	taskdata string
	taskrc   string
	rc       tw.TaskRC

	Context  string
	Filters  tw.Filters
	Active   []string
	Projects map[string]tw.Project
	Contexts map[string]tw.Filters

	logpath string
}

func (m *Model) LoadRc() {
	var err error

	// Load taskrc
	m.taskdata, m.taskrc = actions.Paths()
	m.rc, err = tw.LoadTaskRC(m.taskrc)
	if err != nil {
		log.Fatalln("Could not load taskrc", err)
	}
	log.Println("taskdata", m.taskdata)
	log.Println("taskrc", m.taskrc)

	// Load available contexts with filters from rc
	m.Contexts = m.rc.Contexts

	// Get current context
	context, read_filters, write_filters, err := actions.Context()
	if err != nil {
		log.Fatalln("Could not retrieve context", err)
	}
	m.Context = context
	m.Filters = tw.Filters{
		Read:  read_filters,
		Write: write_filters,
	}

	// Get projects as tree
	m.Projects, err = actions.Projects()
	if err != nil {
		log.Fatalln("Could not retrieve projects", err)
	}

	// Check active
	m.Active, err = actions.Active()
	if err != nil {
		log.Fatalln("Could not retrieve active", err)
	}
}
