package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/manyids2/tasktea/task"
	"github.com/manyids2/tasktea/ui/app"
)

type Styles struct {
	Frame lg.Style
	Focus lg.Style
	Blur  lg.Style
}

func DefaultStyles() Styles {
	return Styles{
		Focus: lg.NewStyle().BorderLeft(true).BorderStyle(lg.NormalBorder()),
		Blur:  lg.NewStyle().BorderLeft(true).BorderStyle(lg.HiddenBorder()),
		Frame: lg.NewStyle().Width(40).Height(30),
	}
}

type Model struct {
	altscreen bool
	quitting  bool
	width     int
	height    int
	focus     int

	styles Styles

	apps []*app.Model
}

func New(projects []task.Project, contexts []task.Context) Model {
	N := 3
	apps := make([]*app.Model, N)
	for i := range projects[:N] {
		a := app.New()
		a.Project = projects[i]
		a.Context = contexts[i]
		a.LoadTasks()
		apps[i] = &a
	}

	return Model{
		apps:   apps,
		styles: DefaultStyles(),
		focus:  0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	if m.quitting {
		return "Bye!\n"
	}
	view := ""
	var style lg.Style
	for i, a := range m.apps {
		switch m.focus == i {
		case true:
			style = m.styles.Focus
		case false:
			style = m.styles.Blur
		}
		view = lg.JoinHorizontal(lg.Top, view,
			style.Render(m.styles.Frame.Render(a.View())))
	}
	return view
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "tab":
			m.focus = (m.focus + 1) % len(m.apps)
			return m, nil
		}
	}
	ma, cmd := m.apps[m.focus].Update(msg)
	m.apps[m.focus] = &ma
	return m, cmd
}
