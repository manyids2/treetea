package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/manyids2/tasktea/task"
	"github.com/manyids2/tasktea/tree"
)

type Styles struct {
	Project lg.Style
	Context lg.Style
	Tasks   lg.Style

	Normal  lg.Style
	Current lg.Style
}

func DefaultStyles() Styles {
	return Styles{
		Project: lg.NewStyle().Bold(true).Padding(1, 2, 0, 2),
		Context: lg.NewStyle().Italic(true).Padding(1, 2, 0, 0),
		Tasks:   lg.NewStyle().Padding(0, 0, 0, 1),

		Normal:  lg.NewStyle().Inline(true),
		Current: lg.NewStyle().Italic(true).Inline(true),
	}
}

type Model struct {
	altscreen bool
	quitting  bool
	width     int
	height    int

	Styles Styles

	Project task.Project
	Context task.Context
	Tasks   *[]tree.Item
	tree    tree.Tree
}

func New() Model {
	m := Model{
		Styles:  DefaultStyles(),
		Project: task.Project{Name: "Inbox", Subprojects: []string{}},
		Context: task.Context{Name: "none", Read: "", Write: ""},
	}
	return m
}

func (m *Model) LoadTasks() {
	tasks, err := task.List(m.Context.Read)
	if err != nil {
		// TODO: Sent to statusbar
		fmt.Println(err)
	} else {
		// HACK: Unfortunately have to copy
		tasks_ := []tree.Item{}
		for _, t := range tasks {
			tasks_ = append(tasks_, tree.Item(t))
		}
		m.Tasks = &tasks_
		m.tree = tree.New(m.Tasks)
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Render() string {
	out := "\n"
	for i, v := range m.tree.Order {
		level := m.tree.Levels[v]
		indent := strings.Repeat("  ", level)
		var style lg.Style
		if i == m.tree.Current() {
			style = m.Styles.Current
		} else {
			style = m.Styles.Normal
		}
		text := fmt.Sprintf("%s %s", indent, ((*m.Tasks)[i]).(task.Task))
		out += style.Render(text) + "\n"
	}
	return out
}

func (m Model) View() string {
	view := ""
	view = lg.JoinHorizontal(lg.Top,
		m.Styles.Project.Render(m.Project.Name),
		m.Styles.Context.Render(m.Context.Name))
	view = lg.JoinVertical(lg.Left, view,
		m.Styles.Tasks.Render(m.Render()))
	return view
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
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
