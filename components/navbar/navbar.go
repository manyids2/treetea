package navbar

import (
	"strings"

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
	StateFilter
)

type Model struct {
	Styles  Styles
	Context string
	Filters []string
	State   State
	Width   int
	Height  int
	Padding string
}

func NewModel(context string, filters []string) Model {
	return Model{
		Styles:  NewStyles(),
		Context: context,
		Filters: filters,
		Width:   24,
		Height:  1,
		Padding: "  ",
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
		m.Padding +
		m.viewIcon() +
		m.viewContext() +
		m.Styles.Icons.Render("  "+ccc.IconUIBar) +
		m.viewFilters()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	// --- RESIZE ---
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
	}
	return m, cmd
}
