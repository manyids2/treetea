package tree

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lgs "github.com/charmbracelet/lipgloss"
	ccc "github.com/manyids2/tasktea/components/theme"
)

func max(a, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

func min(a, b int) int {
	if b >= a {
		return a
	} else {
		return b
	}
}

// Assumed item interface
type Item interface {
	Key() string
	Val() string
	Children() []string
	Extra(string) string
}

// Can later generalize
type Items []Item

// Search array instead of using a map
func (t Items) Get(key string) Item {
	for _, v := range t {
		if key == v.Key() {
			return v
		}
	}
	return nil
}

// Again, search the array
func (t Items) Index(key string) int {
	for k, v := range t {
		if key == v.Key() {
			return k
		}
	}
	return -1
}

// --- Style ---
type Styles struct {
	NormalIcon  lgs.Style
	NormalText  lgs.Style
	CurrentIcon lgs.Style
	CurrentText lgs.Style
	EditIcon    lgs.Style
	EditText    lgs.Style
	ExtraText   lgs.Style
}

func NewStyles() Styles {
	return Styles{
		NormalIcon:  lgs.NewStyle().Foreground(lgs.Color(ccc.ColorMutedBackground)),
		NormalText:  lgs.NewStyle().Foreground(lgs.Color(ccc.ColorMutedForeground)),
		CurrentIcon: lgs.NewStyle().Foreground(lgs.Color(ccc.ColorAlert)).Bold(true),
		CurrentText: lgs.NewStyle().Background(lgs.Color(ccc.ColorEmphBackground)).Foreground(lgs.Color(ccc.ColorForeground)).Italic(true),
		EditIcon:    lgs.NewStyle().Foreground(lgs.Color(ccc.ColorAlert)).Bold(true),
		EditText:    lgs.NewStyle().Background(lgs.Color(ccc.ColorEmphBackground)).Foreground(lgs.Color(ccc.ColorAlert)).Italic(true),
		ExtraText:   lgs.NewStyle().Foreground(lgs.Color(ccc.ColorExtraForeground)),
	}
}

// --- State ---
type State int

const (
	StateHome = iota
	StateEdit
	StateAdd
)

// --- Model ---
type Model struct {
	Items    Items             // Should be satisfied by Task
	Levels   map[string]int    // Needed to track indent
	Parents  map[string]string // Reverse tree
	Children []string          // Top level items
	Order    []string          // Current viewing order
	Extra    string            // Extra info to show for each item

	// State
	State State

	// Appearance
	Width   int
	Height  int
	Padding string
	Styles  Styles

	// Keeping track of the list
	Current int

	// Keeping track of parent when adding
	Parent string

	// Pages
	pages paginator.Model
	input textinput.Model
	err   error
}

func New() Model {
	m := Model{
		Items:    Items{},
		Parents:  map[string]string{},
		Levels:   map[string]int{},
		Children: []string{},
		Order:    []string{},

		Width:   80,
		Height:  20,
		Padding: "  ",
		Styles:  NewStyles(),

		State: StateHome,

		Current: -1,

		pages: paginator.New(),
		input: textinput.New(),
	}
	m.pages.Type = paginator.Dots
	m.pages.PerPage = 10 // statusbar
	m.pages.SetTotalPages(len(m.Items))
	m.input.Prompt = ""
	m.input.Placeholder = ""
	return m
}

func (m *Model) Load(items Items) {
	// Keep ref to old current key, account for first load
	old_key := ""
	if (m.Current >= 0) && (len(m.Order) > m.Current) {
		old_key = m.Order[m.Current]
	}

	// Reset
	m.Items = items
	m.Parents = map[string]string{}
	m.Levels = map[string]int{}
	m.Children = []string{}
	m.Order = []string{}

	// Reverse tree to keep track of parents
	for _, n := range m.Items {
		for _, c := range n.Children() {
			m.Parents[c] = n.Key()
		}
	}

	// Index levels and top level items
	for _, n := range m.Items {
		if m.Parents[n.Key()] == "" {
			m.IndexLevels(n.Key(), 0)
			m.IndexOrder(n.Key())
			m.Children = append(m.Children, n.Key())
		}
	}

	// Really need a get for order as well
	m.Current = 0
	for idx, id := range m.Order {
		if id == old_key {
			m.Current = idx
		}
	}

	// Update paginator
	m.pages.PerPage = m.Height
	m.pages.SetTotalPages(len(m.Items))

	// Go to page of current - may need math
	m.pages.Page = m.Current / m.pages.PerPage
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) CurrentItem() Item {
	id := m.Order[m.Current]
	item := m.Items.Get(id)
	return item
}

func (m *Model) IndexLevels(id string, level int) {
	m.Levels[id] = level
	n := m.Items.Get(id)
	if n == nil {
		return
	}
	for _, c := range n.Children() {
		m.IndexLevels(c, level+1)
	}
}

func (m *Model) IndexOrder(id string) {
	n := m.Items.Get(id)
	if n == nil {
		return
	}
	m.Order = append(m.Order, id)
	for _, c := range n.Children() {
		m.IndexOrder(c)
	}
}

func (m Model) viewIcon(id string) string {
	style := m.Styles.NormalIcon
	icon := ccc.IconUIBar
	if m.Order[m.Current] == id {
		style = m.Styles.CurrentIcon
		icon = ccc.IconUICurrent
	}
	return style.Render("  " + icon)
}

func (m Model) viewPages() string {
	style := m.Styles.NormalText
	return style.Render(fmt.Sprintf("%d items", len(m.Items)))
}

func (m Model) viewItem(id string) string {
	// Get data
	n := m.Items.Get(id)

	// Render
	add_text := ""
	style := m.Styles.NormalText
	text := fmt.Sprintf("%v", n)
	indent := strings.Repeat("  ", m.Levels[id])
	if m.Order[m.Current] == id {
		switch m.State {
		case StateEdit:
			style = m.Styles.EditText
			text = fmt.Sprintf("   %s", m.input.View())
		case StateAdd:
			add_style := m.Styles.EditText
			add_text = fmt.Sprintf("%s %s%s\n",
				m.Padding,
				m.viewIcon(id),
				add_style.Render(fmt.Sprintf("%s    %s", indent, m.input.View())),
			)
		default:
			style = m.Styles.CurrentText
		}
	}

	var content string
	switch m.State {

	case StateEdit:
		content = fmt.Sprintf("%s %s%s\n",
			m.Padding,
			m.viewIcon(id),
			style.Render(fmt.Sprintf("%s %s", indent, text)),
		)

	case StateAdd:
		content = fmt.Sprintf("%s %s%s\n%s",
			m.Padding,
			m.viewIcon(""),
			style.Render(fmt.Sprintf("%s %s", indent, text)),
			add_text,
		)

	default:
		extra := ""
		if m.Extra != "" {
			extra = n.Extra(m.Extra)
		}
		content = fmt.Sprintf("%s %s%s%s\n",
			m.Padding,
			m.viewIcon(id),
			style.Render(fmt.Sprintf("%s %s  ", indent, text)),
			m.Styles.ExtraText.Render(extra))
	}
	return content

}

func (m Model) View() string {
	content := "\n"
	minIdx, maxIdx := m.pages.GetSliceBounds(len(m.Items))
	if maxIdx == minIdx {
		content += fmt.Sprintf("%s No items found.\n\n", m.Padding)
		return content
	}

	for _, v := range m.Order[minIdx:maxIdx] {
		content += m.viewItem(v)
	}

	content += fmt.Sprintf("\n%s %s%s\n",
		m.Padding,
		m.viewIcon(""),
		m.viewPages(),
	)

	if m.pages.TotalPages > 1 {
		content += fmt.Sprintf("\n%s   %s\n", m.Padding, m.pages.View())
	}

	return content
}

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Up         key.Binding
	Down       key.Binding
	Top        key.Binding
	Bottom     key.Binding
	Left       key.Binding
	Right      key.Binding
	Edit       key.Binding
	AddChild   key.Binding
	AddSibling key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "prev page"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "next page"),
	),
	Top: key.NewBinding(
		key.WithKeys("t", "g"),
		key.WithHelp("t/g", "move to top"),
	),
	Bottom: key.NewBinding(
		key.WithKeys("b", "G"),
		key.WithHelp("b/G", "move to bottom"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	AddChild: key.NewBinding(
		key.WithKeys("A"),
		key.WithHelp("A", "add child"),
	),
	AddSibling: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add sibling"),
	),
}

// Returns current input value
type EditChangedMsg string
type EditCancelledMsg string
type errMsg error

func changedEdit(m Model) tea.Cmd {
	return func() tea.Msg {
		return EditChangedMsg(m.input.Value())
	}
}

func cancelEdit(m Model) tea.Cmd {
	return func() tea.Msg {
		task := m.CurrentItem()
		m.input.Placeholder = task.Val()
		m.input.SetValue(task.Val())
		return EditCancelledMsg(m.input.Value()) //
	}
}

func (m *Model) StartEdit() {
	task := m.CurrentItem()
	m.input.Focus()
	m.input.Placeholder = task.Val()
	m.input.SetValue(task.Val())
}

// Returns current input value and parent
type AddChangedMsg []string
type AddCancelledMsg string

func changedAdd(m Model) tea.Cmd {
	key := m.Parent
	m.Parent = ""
	return func() tea.Msg {
		return AddChangedMsg([]string{m.input.Value(), key})
	}
}

func cancelAdd(m Model) tea.Cmd {
	return func() tea.Msg {
		m.input.Placeholder = ""
		m.input.SetValue("")
		m.Parent = ""
		return AddCancelledMsg(m.input.Value()) //
	}
}

func (m *Model) StartAddChild() {
	m.Parent = m.CurrentItem().Key()
	m.input.Placeholder = ""
	m.input.SetValue("")
	m.input.Focus()
}

func (m *Model) StartAddSibling() {
	current := m.CurrentItem().Key()
	m.Parent = m.Parents[current]
	m.input.Placeholder = ""
	m.input.SetValue("")
	m.input.Focus()
}

func (m Model) handleEdit(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.State = StateHome
			return m, changedEdit(m)
		case tea.KeyCtrlC, tea.KeyEsc:
			m.State = StateHome
			return m, cancelEdit(m)
		}
	case errMsg:
		m.err = msg
		return m, nil
	}
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) handleAdd(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.State = StateHome
			return m, changedAdd(m)
		case tea.KeyCtrlC, tea.KeyEsc:
			m.State = StateHome
			return m, cancelAdd(m)
		}
	case errMsg:
		m.err = msg
		return m, nil
	}
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) handleHome(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		minIdx, maxIdx := m.pages.GetSliceBounds(len(m.Items))
		switch {
		// Up
		case key.Matches(msg, keys.Up):
			m.Current = max(m.Current-1, minIdx)

		// Down
		case key.Matches(msg, keys.Down):
			m.Current = min(m.Current+1, maxIdx-1)

		// Top
		case key.Matches(msg, keys.Top):
			m.Current = minIdx

		// Bottom
		case key.Matches(msg, keys.Bottom):
			m.Current = maxIdx - 1

		// Top
		case key.Matches(msg, keys.Left):
			m.pages.PrevPage()
			m.Current = m.pages.PerPage * m.pages.Page

		// Bottom
		case key.Matches(msg, keys.Right):
			m.pages.NextPage()
			m.Current = m.pages.PerPage * m.pages.Page

		// Edit
		case key.Matches(msg, keys.Edit):
			m.State = StateEdit
			m.StartEdit()

		// Add child
		case key.Matches(msg, keys.AddChild):
			m.State = StateAdd
			m.StartAddChild()

		// Add sibling
		case key.Matches(msg, keys.AddSibling):
			m.State = StateAdd
			m.StartAddSibling()
		}
	}
	return m, cmd
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.State {
	case StateHome:
		m, cmd = m.handleHome(msg)
	case StateEdit:
		m, cmd = m.handleEdit(msg)
	case StateAdd:
		m, cmd = m.handleAdd(msg)
	default: // should never occur
	}
	return m, cmd
}
