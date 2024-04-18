package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/manyids2/tasktea/task"
	"github.com/manyids2/tasktea/tree"
	"github.com/manyids2/tasktea/ui/app"
)

type Scene int

const (
	SceneProjects Scene = iota
	SceneContexts
	SceneTags
	SceneTasks
)

type Styles struct {
	Frame lg.Style
	Focus lg.Style
	Blur  lg.Style

	Normal  lg.Style
	Current lg.Style
}

func DefaultStyles() Styles {
	return Styles{
		Focus: lg.NewStyle().BorderLeft(true).BorderStyle(lg.NormalBorder()),
		Blur:  lg.NewStyle().BorderLeft(true).BorderStyle(lg.HiddenBorder()),
		Frame: lg.NewStyle().Width(40).Height(30),

		Normal:  lg.NewStyle().Inline(true),
		Current: lg.NewStyle().Italic(true).Inline(true),
	}
}

type Model struct {
	altscreen bool
	quitting  bool
	width     int
	height    int
	focus     int

	styles Styles

	apps   []*app.Model
	scenes []Scene

	// Common to all apps
	Projects  []task.Project
	projects  tree.Tree
	projects_ *[]tree.Item // Dont know how to avoid copy

	Contexts  []task.Context
	contexts  tree.Tree
	contexts_ *[]tree.Item // Dont know how to avoid copy

	Tags  []task.Tag
	tags  tree.Tree
	tags_ *[]tree.Item // Dont know how to avoid copy
}

func New(projects []task.Project, contexts []task.Context) Model {
	N := 3
	apps := make([]*app.Model, N)
	scenes := make([]Scene, N)
	for i := range projects[:N] {
		a := app.New()
		a.Project = projects[i]
		a.Context = contexts[i]
		a.LoadTasks()
		apps[i] = &a
		scenes[i] = SceneTasks
	}

	m := Model{
		apps:     apps,
		styles:   DefaultStyles(),
		focus:    0,
		scenes:   scenes,
		Projects: projects,
		Contexts: contexts,
		Tags:     []task.Tag{},
	}

	// HACK: Unfortunately have to copy
	projects_ := []tree.Item{}
	for _, t := range projects {
		projects_ = append(projects_, tree.Item(t))
	}
	m.projects_ = &projects_
	m.projects = tree.New(&projects_)

	// HACK: Unfortunately have to copy
	contexts_ := []tree.Item{}
	for _, t := range contexts {
		contexts_ = append(contexts_, tree.Item(t))
	}
	m.contexts_ = &contexts_
	m.contexts = tree.New(&contexts_)

	// HACK: Unfortunately have to copy
	tags_ := []tree.Item{}
	for _, t := range m.Tags {
		tags_ = append(tags_, tree.Item(t))
	}
	m.tags_ = &tags_
	m.tags = tree.New(&tags_)

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) RenderProjects() string {
	out := "\n"
	for i, v := range m.projects.Order {
		level := m.projects.Levels[v]
		indent := strings.Repeat("  ", level)
		var style lg.Style
		if i == m.projects.Current() {
			style = m.styles.Current
		} else {
			style = m.styles.Normal
		}
		text := fmt.Sprintf("%s %s", indent, ((*m.projects_)[i]).(task.Project).Name)
		out += style.Render(text) + "\n"
	}
	return out
}

func (m Model) RenderContexts() string {
	out := "\n"
	for i, v := range m.contexts.Order {
		level := m.contexts.Levels[v]
		indent := strings.Repeat("  ", level)
		var style lg.Style
		if i == m.contexts.Current() {
			style = m.styles.Current
		} else {
			style = m.styles.Normal
		}
		text := fmt.Sprintf("%s %s", indent, ((*m.contexts_)[i]).(task.Context).Name)
		out += style.Render(text) + "\n"
	}
	return out
}

func (m Model) RenderTags() string {
	out := "\n"
	for i, v := range m.tags.Order {
		level := m.tags.Levels[v]
		indent := strings.Repeat("  ", level)
		var style lg.Style
		if i == m.tags.Current() {
			style = m.styles.Current
		} else {
			style = m.styles.Normal
		}
		text := fmt.Sprintf("%s %s", indent, ((*m.tags_)[i]).(task.Tag))
		out += style.Render(text) + "\n"
	}
	return out
}

func (m Model) View() string {
	if m.quitting {
		return "Bye!\n"
	}
	view := ""
	var style lg.Style
	for i, a := range m.apps {
		switch m.focus == i {
		case true:
			style = m.styles.Focus
		case false:
			style = m.styles.Blur
		}

		var text string
		switch m.scenes[i] {
		case SceneTasks:
			text = a.View()
		case SceneProjects:
			text = m.RenderProjects()
		case SceneContexts:
			text = m.RenderContexts()
		case SceneTags:
			text = m.RenderTags()
		}

		text = style.Render(m.styles.Frame.Render(text))
		view = lg.JoinHorizontal(lg.Top, view, text)
	}
	return view
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			m.scenes[m.focus] = SceneTasks
			return m, nil
		case "p":
			m.scenes[m.focus] = SceneProjects
			return m, nil
		case "c":
			m.scenes[m.focus] = SceneContexts
			return m, nil
		case "t":
			m.scenes[m.focus] = SceneTags
			return m, nil
		case "tab":
			m.focus = (m.focus + 1) % len(m.apps)
			return m, nil
		}
	}
	ma, cmd := m.apps[m.focus].Update(msg)
	m.apps[m.focus] = &ma
	return m, cmd
}
