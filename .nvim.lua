local ls = require("luasnip")
local s = ls.snippet
local t = ls.text_node

ls.add_snippets("all", {
	s("trig", {
		t("Worked!!"),
	}),
})

ls.add_snippets("go", {
	s("Tea", {
		t({
			"package item",
			"import (",
			'	tea "github.com/charmbracelet/bubbletea"',
			")",
			"",
			"type Model struct {",
			"}",
			"",
			"func New() Model {",
			"	return Model{}",
			"}",
			"",
			"func (m Model) Init() tea.Cmd {",
			"	return nil",
			"}",
			"",
			"func (m Model) View() string {",
			'	return ""',
			"}",
			"",
			"func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {",
			"	return m, nil",
			"}",
		}),
	}),
})
