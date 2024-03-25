package statusbar

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Message string
}

func New() Model {
	return Model{Message: ""}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	return m.Message
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
