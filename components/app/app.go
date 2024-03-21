package navbar

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lgs "github.com/charmbracelet/lipgloss"
	"github.com/manyids2/tasktea/components/navbar"
	"github.com/manyids2/tasktea/components/statusbar"
	ccc "github.com/manyids2/tasktea/components/theme"
)

type Styles struct {
	Normal lgs.Style
}

func NewStyles() Styles {
	return Styles{
		Normal: lgs.NewStyle().Foreground(lgs.Color(ccc.ColorMutedForeground)),
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

	nav    navbar.Model
	status statusbar.Model
	help   help.Model
	keys   keyMap

	quitting bool
}

type keyMap struct {
	Quit   key.Binding
	Filter key.Binding
	Help   key.Binding
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q/esc", "quit"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Filter, k.Help, k.Quit},
	}
}

func NewModel(context string, filters []string) Model {
	return Model{
		Context: context,
		Filters: filters,
		State:   StateHome,
		Width:   80,
		Height:  1,
		Padding: "  ",
		Styles:  NewStyles(),
		nav:     navbar.NewModel(context, filters),
		status:  statusbar.NewModel(),
		help:    help.New(),
		keys:    keys,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	return m.nav.View() +
		"\n\n" +
		m.help.View(m.keys)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// --- RESIZE ---
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
	}

	// --- Handle keys based on current state ---
	switch m.State {

	// --- Home ---
	case StateHome:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			// Quit
			case key.Matches(msg, keys.Quit):
				m.quitting = true
				return m, tea.Quit

			// Help
			case key.Matches(msg, m.keys.Help):
				m.help.ShowAll = !m.help.ShowAll

			// Edit filter
			case key.Matches(msg, keys.Filter):
				m.State = StateFilter
			}
		}

	// --- Edit filter ---
	case StateFilter:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			// Quit
			case key.Matches(msg, keys.Quit):
				m.quitting = true
				return m, tea.Quit

			case key.Matches(msg, keys.Filter):
				m.State = StateHome
			}
		}

	// should never occur
	default:
	}

	return m, cmd
}
