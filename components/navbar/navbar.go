package navbar

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lgs "github.com/charmbracelet/lipgloss"
	ccc "github.com/manyids2/tasktea/components/theme"
)

type Styles struct {
	Filters lgs.Style
	Context lgs.Style
	Icons   lgs.Style
}

func NewStyles() Styles {
	return Styles{
		Filters: lgs.NewStyle().Foreground(lgs.Color(ccc.ColorMutedForeground)),
		Context: lgs.NewStyle().Foreground(lgs.Color(ccc.ColorForeground)).Bold(true),
		Icons:   lgs.NewStyle().Foreground(lgs.Color(ccc.ColorForeground)).Bold(true),
	}
}

type State int

const (
	StateHome = iota
)

type Model struct {
	Styles   Styles
	Context  string
	Filters  []string
	state    State
	width    int
	height   int
	padding  string
	quitting bool
}

func NewModel(context string, filters []string) Model {
	return Model{
		Styles:  NewStyles(),
		Context: context,
		Filters: filters,
		width:   24,
		height:  1,
		padding: "  ",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) viewIcon() string {
	return m.Styles.Icons.Render(ccc.IconUIFilter + " " + ccc.IconUIBar)
}

func (m Model) viewContext() string {
	return m.Styles.Context.Render(m.Context)
}

func (m Model) viewFilters() string {
	return m.Styles.Filters.Render(strings.Join(m.Filters, " "))
}

func (m Model) View() string {
	//
	//   ▍ project:jobsearch -init
	//
	return "\n" +
		m.padding +
		m.viewIcon() +
		m.viewContext() +
		m.Styles.Icons.Render("  "+ccc.IconUIBar) +
		m.viewFilters()
}

type keyMap struct {
	Quit key.Binding
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q/esc", "quit"),
	),
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// --- RESIZE ---
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}

	// --- Handle keys based on current state ---
	switch m.state {
	case StateHome:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			// Quit
			case key.Matches(msg, keys.Quit):
				m.quitting = true
				return m, tea.Quit
			}
		}

	// should never occur
	default:
	}

	return m, cmd
}
