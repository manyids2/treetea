package tree

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	tk "github.com/manyids2/tasktea/task"
)

type Items []tk.Task

// Search array instead of using a map
func (t Items) Get(key string) tk.Task {
	for _, v := range t {
		if key == v.UUID {
			return v
		}
	}
	return tk.Task{}
}

// --- State ---
type State int

const (
	StateHome = iota
	StateNew
	StateEdit
)

// --- Model ---
type Model struct {
	Items   Items             // Should be satisfied by Task
	Levels  map[string]int    // Needed to track indent
	Parents map[string]string // Reverse tree

	State State

	Width   int
	Height  int
	Padding string
}

func New(items []tk.Task) Model {
	m := Model{
		Items:   items,
		Parents: map[string]string{},
		Levels:  map[string]int{},

		Width:   64,
		Height:  1,
		Padding: "  ",

		State: StateHome,
	}

	// Reverse tree to keep track of parents
	for _, item := range m.Items {
		for _, c := range item.Depends {
			m.Parents[c] = item.UUID
		}
	}

	// Root items at level 0
	for _, item := range m.Items {
		if m.Parents[item.UUID] == "" {
			m.IndexLevels(item.UUID, 0)
		}
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) IndexLevels(id string, level int) {
	m.Levels[id] = level
	n := m.Items.Get(id)
	for _, c := range n.Depends {
		m.IndexLevels(c, level+1)
	}
}

func (m Model) viewTree(n tk.Task, content string) string {
	content += fmt.Sprintf("%s %s\n", strings.Repeat("  ", m.Levels[n.UUID]), n.UUID)
	for _, v := range n.Depends {
		c := m.Items.Get(v) // We are certain to find the uuid
		content = m.viewTree(c, content)
	}
	return content
}

func (m Model) View() string {
	content := ""

	// content += "Levels\n"
	// for k, v := range m.Levels {
	// 	content += fmt.Sprintf("%s: %d\n", k, v)
	// }

	// content += "Parents\n"
	// for k, v := range m.Parents {
	// 	content += fmt.Sprintf("%s: %s\n", k, v)
	// }

	// content += "Items\n"
	// for _, c := range m.Items {
	// 	content += c.UUID + "\n"
	// }

	content += "\n"
	for _, c := range m.Items {
		if (m.Parents[c.UUID] == "") && c.UUID != "" {
			content = m.viewTree(c, content)
		}
	}

	return content
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
