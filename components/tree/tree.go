package tree

import (
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

// Assumed item interface
type Item interface {
	key() string
	children() []string
	desc(string) string
}

// Can later generalize
type Items []Item

// Search array instead of using a map: check i >=0 for available
func (t Items) Get(key string) (i int, v Item) {
	for i, v := range t {
		if key == v.key() {
			return i, v
		}
	}
	return -1, nil
}

type Model struct {
	Width  int
	Height int

	Items    Items             // Should be satisfied by Task
	Levels   map[string]int    // Needed to track indent
	Parents  map[string]string // Reverse tree
	Children []string          // Top level items
	Order    []string          // Current viewing order
	Extra    string            // Extra info to show for each item
	Active   string            // Currently active task
	Selected []string          // Selected tasks

	Current int
	Parent  string

	pages paginator.Model
	frame lg.Style
	err   error
}

func (m *Model) SetFrame(width, height int) {
	m.Width, m.Height = width, height
	m.frame = lg.NewStyle().Height(height).Width(width)
}

func New() (m Model) {
	m = Model{
		Items:    Items{},
		Parents:  map[string]string{},
		Levels:   map[string]int{},
		Children: []string{},
		Order:    []string{},
		Selected: []string{},
		Width:    80,
		Height:   3,
		Current:  -1,
		pages:    paginator.New(),
	}
	m.SetFrame(m.Width, m.Height)
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Reset(items Items) {
	m.Items = items
	m.Parents = map[string]string{}
	m.Levels = map[string]int{}
	m.Children = []string{}
	m.Order = []string{}
	m.Selected = []string{}
	m.pages.SetTotalPages(len(m.Items))
}

func (m *Model) ResetCurrent() {
	if len(m.Items) == 0 {
		m.Current = -1
		m.pages.Page = 0
		return
	}
	m.Current = 0
	m.pages.Page = m.Current / m.pages.PerPage
}

func (m *Model) IndexLevels(id string, level int) {
	m.Levels[id] = level
	_, n := m.Items.Get(id)
	if n == nil {
		return
	}
	for _, c := range n.children() {
		m.IndexLevels(c, level+1)
	}
}

func (m *Model) IndexOrder(id string) {
	_, n := m.Items.Get(id)
	if n == nil {
		return
	}
	m.Order = append(m.Order, id)
	for _, c := range n.children() {
		m.IndexOrder(c)
	}
}

func (m *Model) LoadTree(items Items) {
	m.Reset(items)
	m.ResetCurrent()

	// Reverse tree to keep track of parents
	for _, v := range m.Items {
		for _, c := range v.children() {
			m.Parents[c] = v.key()
		}
	}

	// Index levels and top level items
	for _, v := range m.Items {
		k := v.key()
		if m.Parents[k] == "" {
			m.IndexLevels(k, 0)
			m.IndexOrder(k)
			m.Children = append(m.Children, k)
		}
	}
}

func (m *Model) LoadList(items Items) {
	m.Reset(items)
	m.ResetCurrent()

	// Reverse tree to keep track of parents
	for _, v := range m.Items {
		m.Parents[v.key()] = ""
	}

	// Index levels and top level items
	for _, v := range m.Items {
		m.Order = append(m.Order, v.key())
	}
	m.Children = m.Order
}

func (m Model) View() string {
	return m.frame.Render("tree")
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
