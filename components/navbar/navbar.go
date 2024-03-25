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
}

func New() Model {
	return Model{
		Width:  80,
		Height: 3,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	var style = lg.NewStyle().
		Height(m.Height).
		Width(m.Width)
	return style.Render(fmt.Sprintf("%s | %s", m.Title, m.Description))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
