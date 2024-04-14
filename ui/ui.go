package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	altscreen bool
	quitting  bool
	width     int
	height    int
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	var mode string
	if m.altscreen {
		mode = " altscreen mode "
	} else {
		mode = " inline mode "
	}
	mode = fmt.Sprintf("%s [ %d, %d ]", mode, m.width, m.height)

	return mode
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
