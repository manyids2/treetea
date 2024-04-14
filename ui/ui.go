package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/manyids2/tasktea/ui/app"
)

type Model struct {
	altscreen bool
	quitting  bool
	width     int
	height    int
	focus     int

	apps []*app.Model
}

func New() Model {
	a := app.New()
	a.Project.Name = "Project a"
	a.Focus = true

	b := app.New()
	b.Project.Name = "Project b"

	return Model{
		apps:  []*app.Model{&a, &b},
		focus: 0,
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
	for _, a := range m.apps {
		view = lg.JoinHorizontal(lg.Top, view, a.View())
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
			m.apps[m.focus].Focus = false
			m.focus = (m.focus + 1) % len(m.apps)
			m.apps[m.focus].Focus = true
			return m, nil
		}
	}
	ma, cmd := m.apps[m.focus].Update(msg)
	m.apps[m.focus] = &ma
	return m, cmd
}
