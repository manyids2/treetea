package app

import "fmt"

func (m Model) View() (out string) {
	if m.quitting {
		return ""
	}
	return m.Context + "\n" +
		m.Filters.Read + "\n" +
		fmt.Sprintf("%d tasks", len(m.Tasks))
}
