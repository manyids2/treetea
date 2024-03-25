package navbar

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type Model struct {
	Width       int
	Height      int
	Title       string
	Description string
	frame       lg.Style
}

func New() (m Model) {
	m = Model{
		Width:  80,
		Height: 3,
	}
	m.SetFrame(m.Width, m.Height)
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) SetFrame(width, height int) {
	m.Width, m.Height = width, height
	m.frame = lg.NewStyle().Height(height).Width(width)
}

func (m Model) View() string {
	return m.frame.Render(fmt.Sprintf("\n| %s | %s", m.Title, m.Description))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
