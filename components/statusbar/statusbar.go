package statusbar

import (
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type Model struct {
	Width   int
	Height  int
	Message string
}

func New() Model {
	return Model{
		Width:   80,
		Height:  2,
		Message: "statusbar",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	var style = lg.NewStyle().
		Height(m.Height).
		Width(m.Width)
	return style.Render(m.Message)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
