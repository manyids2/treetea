package app

import (
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/manyids2/tasktea/task"
)

type Styles struct {
	Project lg.Style
	Context lg.Style
}

func DefaultStyles() Styles {
	return Styles{
		Project: lg.NewStyle().Bold(true).Padding(1, 2),
		Context: lg.NewStyle().Italic(true).Padding(1, 2, 1, 0),
	}
}

type Model struct {
	altscreen bool
	quitting  bool
	width     int
	height    int

	Styles Styles

	Project task.Project
	Context task.Context
	Tasks   []task.Task
}

func New() Model {
	return Model{
		Styles:  DefaultStyles(),
		Project: task.Project{Name: "Project"},
		Context: task.Context{Name: "Context"},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	view := ""

	view = lg.JoinHorizontal(lg.Top,
		m.Styles.Project.Render(m.Project.Name),
		m.Styles.Context.Render(m.Context.Name))
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
