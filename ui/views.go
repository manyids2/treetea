package ui

import (
	"fmt"
	"strings"

	tk "github.com/manyids2/tasktea/task"
)

const bar = "▍ "

// filtersView Render filters, input if editing
func (m Model) filtersView() string {
	content := m.styles.glass.Render("   " + bar)
	if m.state == StateFilter {
		content += m.styles.edit.Render(m.inputFilters.View())
	} else {
		content += m.styles.filters.Render(strings.Join(m.filters, " "))
	}
	return content + "\n"
}

// contentView Render indented task list, input if editing
func (m Model) statusView() string {
	padding := "    "
	return fmt.Sprintf("\n%s%s%d tasks\n", padding, bar, len(m.items))
}

// contentView Render indented task list, input if editing
func (m Model) contentView() string {
	padding := "    "
	content := ""
	// Loop over items to preserve order
	for i, uuid := range m.items {
		// Get task
		task := m.tasks[uuid]
		inner := tk.View(m.levels[task.UUID], task.Status, task.Description)

		// Style for task
		prefix := "  "
		style := m.styles.normal

		// Modify for current
		if i == m.current {
			prefix = m.styles.current.Render(bar)
			style = m.styles.italic
		}

		// Final task line
		content += padding + prefix + style.Render(inner) + "\n"
	}
	return content
}

// View Render complete UI
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	return "\n" + m.filtersView() + // Filter
		"\n" + m.contentView() + // Tasks
		m.statusView() + // Status
		"\n" + m.help.View(m.keys) // Help
}
