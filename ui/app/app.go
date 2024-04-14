package app

import (
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type Styles struct {
	Frame   lg.Style
	Focus   lg.Style
	Blur    lg.Style
	Project lg.Style
	Context lg.Style
}

func DefaultStyles() Styles {
	return Styles{
		Focus:   lg.NewStyle().BorderLeft(true).BorderStyle(lg.BlockBorder()),
		Blur:    lg.NewStyle().BorderLeft(true).BorderStyle(lg.NormalBorder()),
		Project: lg.NewStyle().Bold(true).Padding(2, 2),
		Context: lg.NewStyle().Italic(true).Padding(2, 2),
		Frame:   lg.NewStyle().Width(40).Height(30),
	}
}

type Project struct {
	Name     string
	Children []string
}

type Context struct {
	Name  string
	Read  string // Filters
	Write string
}

type Model struct {
	altscreen bool
	quitting  bool
	width     int
	height    int

	Styles Styles
	Focus  bool

	Project Project
	Context Context
	Tags    []string

	Tasks []string
}

func New() Model {
	return Model{
		Styles:  DefaultStyles(),
		Project: Project{Name: "Project"},
		Context: Context{Name: "Context"},
		Tags:    []string{},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	view := ""

	var style lg.Style
	switch m.Focus {
	case true:
		style = m.Styles.Focus
	case false:
		style = m.Styles.Blur
	}
	view = style.Render(lg.JoinHorizontal(lg.Top,
		m.Styles.Project.Render(m.Project.Name),
		m.Styles.Context.Render(m.Context.Name)))
	view = m.Styles.Frame.Render(view)
	return view
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case " ":
			var cmd tea.Cmd
			if m.altscreen {
				cmd = tea.ExitAltScreen
			} else {
				cmd = tea.EnterAltScreen
			}
			m.altscreen = !m.altscreen
			return m, cmd
		}
	}
	return m, nil
}
