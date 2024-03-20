package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// handleHome Navigation, help and quit in home state
func (m Model) handleHome(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Keys
	case tea.KeyMsg:
		switch {
		// Help
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		// Quit
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit

		// Up
		case key.Matches(msg, m.keys.Up):
			m.current = max(m.current-1, 0)

		// Down
		case key.Matches(msg, m.keys.Down):
			m.current = min(m.current+1, len(m.items)-1)

		// Top
		case key.Matches(msg, m.keys.Top):
			m.current = 0

		// Bottom
		case key.Matches(msg, m.keys.Bottom):
			m.current = len(m.items) - 1
		}
	}

	return m, nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// NOTE: No global maps, as we have to deal with input
	var cmd tea.Cmd

	// --- RESIZE ---
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	}

	// --- Handle keys based on current state ---
	switch m.state {
	// Home state - navigation
	case StateHome:
		m, cmd = m.handleHome(msg)
	// should never occur
	default:
	}

	return m, cmd
}
