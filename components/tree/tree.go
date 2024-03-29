package tree

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

// Assumed item interface
type Item interface {
	Key() string
	Children() []string
	Desc(string) string
}

// Can later generalize
type Items []Item

// Search array instead of using a map: check i >=0 for available
func (t Items) Get(key string) (i int, v Item) {
	for i, v := range t {
		if key == v.Key() {
			return i, v
		}
	}
	return -1, nil
}

type Model struct {
	Name   string
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

	// For statusbar
	Singular string
	Plural   string
}

func (m *Model) SetFrame(width, height int) {
	m.Width, m.Height = width, height
	m.frame = lg.NewStyle().Height(height).Width(width)
}

const (
	IconActivePage   = "🬋║🬋"
	IconInactivePage = "🬋🬋🬋"
)

func New(name, singular, plural string) (m Model) {
	m = Model{
		Name:     name,
		Singular: singular,
		Plural:   plural,
		Items:    Items{},
		Parents:  map[string]string{},
		Levels:   map[string]int{},
		Children: []string{},
		Order:    []string{},
		Selected: []string{},
		Width:    80,
		Height:   20,
		Current:  -1,
		pages:    paginator.New(),
	}
	m.SetFrame(m.Width, m.Height)
	m.pages.Type = paginator.Dots
	m.pages.ActiveDot = IconActivePage
	m.pages.InactiveDot = IconInactivePage
	m.pages.PerPage = m.Height
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
	m.pages.PerPage = m.Height
	m.pages.SetTotalPages(len(m.Items))
}

func (m *Model) ResetCurrent() {
	if len(m.Items) == 0 {
		m.Current = -1
		m.pages.Page = 0
		return
	}
	m.Current = 0
	m.pages.PerPage = m.Height
	if m.pages.PerPage <= 0 {
		m.pages.PerPage = 20
	}
	m.pages.Page = m.Current / m.pages.PerPage
}

func (m *Model) IndexLevels(id string, level int) {
	m.Levels[id] = level
	_, n := m.Items.Get(id)
	if n == nil {
		return
	}
	for _, c := range n.Children() {
		m.IndexLevels(c, level+1)
	}
}

func (m *Model) IndexOrder(id string) {
	_, n := m.Items.Get(id)
	if n == nil {
		return
	}
	m.Order = append(m.Order, id)
	for _, c := range n.Children() {
		m.IndexOrder(c)
	}
}

func (m *Model) LoadTree(items Items) {
	m.Reset(items)
	m.ResetCurrent()

	// Reverse tree to keep track of parents
	for _, v := range m.Items {
		for _, c := range v.Children() {
			m.Parents[c] = v.Key()
		}
	}

	// Index levels and top level items
	for _, v := range m.Items {
		k := v.Key()
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
		m.Parents[v.Key()] = ""
	}

	// Index levels and top level items
	for _, v := range m.Items {
		m.Order = append(m.Order, v.Key())
	}
	m.Children = m.Order
}

func (m Model) View() string {
	content := ""
	minIdx, maxIdx := m.pages.GetSliceBounds(len(m.Items))

	// Render tree
	for _, id := range m.Order[minIdx:maxIdx] {
		// Get data
		_, n := m.Items.Get(id)
		text := fmt.Sprintf("%s", n.Desc("description"))
		indent := strings.Repeat("  ", m.Levels[id])
		content += fmt.Sprintf("%s%s\n", indent, text)
	}

	// Pad height
	// TODO: Use frame instead
	for i := maxIdx - minIdx; i < m.pages.PerPage; i++ {
		content += "\n"
	}

	// Render status
	content += "\n"
	status := ""
	if len(m.Items) == 0 {
		status = fmt.Sprintf("No %s found", m.Plural)
	} else if len(m.Items) == 1 {
		status = fmt.Sprintf("%d %s", len(m.Items), m.Singular)
	} else {
		status = fmt.Sprintf("%d %s", len(m.Items), m.Plural)
	}
	content += status

	// Render pages
	if m.pages.TotalPages > 1 {
		pages := fmt.Sprintf("%d / %d  %s", m.pages.Page+1, m.pages.TotalPages, m.pages.View())
		padding := strings.Repeat(" ", m.Width-len([]rune(pages))-len([]rune(status)))
		content += fmt.Sprintf("%s%s", padding, pages)
	}
	content += "\n"

	return m.frame.Render(content)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
