package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/manyids2/tasktea/components/layout"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Non-key messages ( event listeners )
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case errMsg:
		m.err = msg
		return m, nil

	case layout.ReloadRcMsg:
		// Refresh contexts
		m.LoadRc()
		m.layout.LoadContexts(m.Contexts)
		return m, nil

	case layout.CloseMsg:
		m.quitting = true
		return m, tea.Quit
	}

	// Defer to layout
	m.layout, cmd = m.layout.Update(msg)
	return m, cmd
}
