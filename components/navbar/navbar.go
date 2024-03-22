package navbar

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lgs "github.com/charmbracelet/lipgloss"
	ccc "github.com/manyids2/tasktea/components/theme"
)

type Styles struct {
	Filters lgs.Style
	Editing lgs.Style
	Context lgs.Style
	Icons   lgs.Style
}

func NewStyles() Styles {
	return Styles{
		Filters: lgs.NewStyle().Foreground(lgs.Color(ccc.ColorMutedForeground)),
		Editing: lgs.NewStyle().Foreground(lgs.Color(ccc.ColorAlert)).Italic(true),
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
	Context string
	Filters []string

	State State

	Width   int
	Height  int
	Padding string
	Styles  Styles

	input textinput.Model
	err   error
}

func New(context string, filters []string) Model {
	m := Model{
		Context: context,
		Filters: filters,

		Width:   64,
		Height:  1,
		Padding: "  ",
		Styles:  NewStyles(),

		State: StateHome,

		input: textinput.New(),
	}
	m.input.Prompt = ""
	m.input.Placeholder = ""
	m.input.SetValue(strings.Join(filters, " "))
	return m
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
	switch m.State {
	case StateHome:
		return m.Styles.Filters.Render(strings.Join(m.Filters, " "))
	case StateFilter:
		return m.Styles.Editing.Render(m.input.View())
	}
	return ""
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

// Returns current input value
type ChangedMsg []string
type CancelledMsg []string
type errMsg error

func updateFilter(m Model) tea.Cmd {
	return func() tea.Msg {
		return ChangedMsg(m.Filters)
	}
}

func cancelFilter(m Model) tea.Cmd {
	return func() tea.Msg {
		return CancelledMsg(m.Filters)
	}
}

func (m *Model) StartFilter() {
	m.State = StateFilter
	m.input.Focus()
}

func (m Model) handleFilter(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.State = StateHome
			m.Filters = strings.Split(m.input.Value(), " ")
			return m, updateFilter(m)
		case tea.KeyCtrlC, tea.KeyEsc:
			m.State = StateHome
			return m, cancelFilter(m)
		}
	case errMsg:
		m.err = msg
		return m, nil
	}
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	// --- RESIZE ---
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
	}

	switch m.State {
	case StateHome: // Nothing to do
	case StateFilter: // Hand off to input
		m, cmd = m.handleFilter(msg)
	default: // Should not occur
	}

	return m, cmd
}
