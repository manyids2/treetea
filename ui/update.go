package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// handleNav Home navigation
func (m Model) handleHome(msg tea.Msg) (Model, tea.Cmd) {
	maxIdx := len(m.items) - 1
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

		// Cancel
		case key.Matches(msg, m.keys.Cancel):
			m.quitting = true
			return m, tea.Quit

		// Up
		case key.Matches(msg, m.keys.Up):
			m.current = max(m.current-1, 0)

		// Down
		case key.Matches(msg, m.keys.Down):
			m.current = min(m.current+1, maxIdx)

		// Top
		case key.Matches(msg, m.keys.Top):
			m.current = 0

		// Bottom
		case key.Matches(msg, m.keys.Bottom):
			m.current = maxIdx - 1

		// History
		case key.Matches(msg, m.keys.History):
			// TODO: Handle no history
			m.keys.Quit.SetHelp("q", "quit")
			m.state = StateHistory
			m.current = 0
		}
	}

	return m, nil
}

// handleHistory History navigation
func (m Model) handleHistory(msg tea.Msg) (Model, tea.Cmd) {
	maxIdx := len(m.history) - 1
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
			m.current = min(m.current+1, maxIdx)

		// Top
		case key.Matches(msg, m.keys.Top):
			m.current = 0

		// Bottom
		case key.Matches(msg, m.keys.Bottom):
			m.current = maxIdx - 1

		// Back home
		case key.Matches(msg, m.keys.History):
			m.keys.Quit.SetHelp("q/esc", "quit")
			m.state = StateHome

		// Accept
		case key.Matches(msg, m.keys.Accept):
			m.filters = strings.Split(m.history[m.current], " ")
			m.state = StateHome
			m.Reload()

		// Cancel
		case key.Matches(msg, m.keys.Cancel):
			m.keys.Quit.SetHelp("q/esc", "quit")
			m.state = StateHome
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
	case StateHome:
		m, cmd = m.handleHome(msg)
	case StateHistory:
		m, cmd = m.handleHistory(msg)

	// should never occur
	default:
	}

	return m, cmd
}
