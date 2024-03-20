package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	tk "github.com/manyids2/tasktea/task"
)

type State int

const (
	StateHome State = iota
	StateView
	StateAdd
	StateEdit
	StateFilter
)

const CHAR_LIMIT = 128

type Model struct {
	state    State              // App State
	filters  []string           // Filters for `task export`
	items    []string           // List of items, order preserved
	tasks    map[string]tk.Task // Task by UUID
	parents  map[string]string  // Parent by UUID
	levels   map[string]int     // Level by UUID
	current  int                // Index of current item
	new_idx  int                // For help with addition, may not be needed
	keys     keyMap             // Global keymaps
	styles   styles             // Lipgloss styles
	err      error              // For error msg
	debug    string             // Debug
	quitting bool               // Mandatory

	// Subcomponents
	inputTask    textinput.Model // Add, edit task
	inputFilters textinput.Model // Edit filters
	help         help.Model      // Short, long help
}

// NewModel Initialize app state, get tasks based on filters
func NewModel(filters []string) Model {
	// Input for task
	inputTask := textinput.New()
	inputTask.Prompt = ""
	inputTask.Placeholder = ""
	inputTask.CharLimit = CHAR_LIMIT
	inputTask.Width = CHAR_LIMIT

	// Input for filters
	inputFilters := textinput.New()
	inputFilters.Prompt = ""
	inputFilters.Placeholder = ""
	inputFilters.CharLimit = CHAR_LIMIT
	inputFilters.Width = CHAR_LIMIT

	// Initialize model
	m := Model{
		state:        StateHome,
		filters:      filters,
		inputTask:    inputTask,
		inputFilters: inputFilters,
		current:      0,
		keys:         keys,
		help:         help.New(),
		styles:       NewStyles(),
		new_idx:      0,
	}

	// Load tasks, parents, levels, items
	m.Reload()
	return m
}

// Init Nop
func (m Model) Init() tea.Cmd {
	return nil
}

// index_levels Recursive function to index levels
func index_levels(task tk.Task, level int, items []string, levels map[string]int, parents map[string]string, tasks map[string]tk.Task) []string {
	levels[task.UUID] = level
	if task.Description != "" {
		items = append(items, task.UUID)
	}
	for _, c := range task.Depends {
		if parents[c] == task.UUID {
			items = index_levels(tasks[c], level+1, items, levels, parents, tasks)
		}
	}
	return items
}

// Reload Reload items by calling `task export [filters]` again
func (m *Model) Reload() {
	list, err := tk.List(m.filters)
	if err != nil {
		// For now just silent return nothing, mark model error
		m.tasks = map[string]tk.Task{}
		m.parents = map[string]string{}
		m.levels = map[string]int{}
		m.items = []string{}
		m.err = err
		return
	}
	tasks := make(map[string]tk.Task, len(list))

	// Index parents
	parents := make(map[string]string, len(list))
	for _, task := range list {
		tasks[task.UUID] = task
		for _, d := range task.Depends {
			parents[d] = task.UUID
		}
	}

	// Index levels
	items := []string{}
	levels := make(map[string]int, len(list))
	for _, v := range list {
		if parents[v.UUID] == "" {
			items = index_levels(v, 0, items, levels, parents, tasks)
		}
	}

	// Update model
	m.tasks = tasks
	m.parents = parents
	m.levels = levels
	m.items = items
	m.err = nil
}
