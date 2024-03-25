package layout

import (
	tea "github.com/charmbracelet/bubbletea"

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

func (m Model) View() string {
	if !m.ready {
		return ""
	}
	return m.nav.View() + "\n" + m.status.View()
}

type errMsg error

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Non-key messages ( event listeners )
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.ready = true

	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, cmd
}
