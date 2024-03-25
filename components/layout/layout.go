package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"

	nv "github.com/manyids2/tasktea/components/navbar"
	st "github.com/manyids2/tasktea/components/statusbar"
	tr "github.com/manyids2/tasktea/components/tree"
)

type ViewState int

const (
	ViewHome = iota
)

type Model struct {
	// Essential
	Width  int
	Height int
	State  ViewState

	// Information
	nav    nv.Model
	status st.Model

	// Trees
	tasks    tr.Model
	contexts tr.Model
	projects tr.Model
	tags     tr.Model
	history  tr.Model

	// Error handling
	ready bool
	err   error
}

func New() Model {
	return Model{
		Width:  80,
		Height: 24,
		nav:    nv.New(),
		status: st.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Layout() {
	m.nav.Width = m.Width
	m.nav.Height = 3

	m.status.Width = m.Width
	m.status.Height = 2
}

func (m Model) View() string {
	var style = lg.NewStyle().
		Background(lg.Color("#dddddd")).
		Height(m.Height).
		Width(m.Width)
	if !m.ready {
		return "not ready"
	}
	return style.Render(m.nav.View() + "\n" + m.status.View())
}

type errMsg error

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	// Non-key messages ( event listeners )
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Layout()
		m.ready = true
		return m, nil

	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, cmd
}
