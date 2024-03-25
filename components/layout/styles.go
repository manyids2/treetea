package layout

import (
	"github.com/charmbracelet/lipgloss"
)

type styles struct {
	glass    lipgloss.Style
	filters  lipgloss.Style
	current  lipgloss.Style
	normal   lipgloss.Style
	italic   lipgloss.Style
	bold     lipgloss.Style
	selected lipgloss.Style
	edit     lipgloss.Style
}

func NewStyles() styles {
	return styles{
		glass:    lipgloss.NewStyle().Foreground(lipgloss.Color("#FFddB7")).Bold(true),
		filters:  lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")).Bold(true),
		current:  lipgloss.NewStyle().Foreground(lipgloss.Color("#00aa00")),
		normal:   lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")),
		italic:   lipgloss.NewStyle().Italic(true),
		bold:     lipgloss.NewStyle().Bold(true),
		selected: lipgloss.NewStyle().Foreground(lipgloss.Color("#aa0000")),
		edit:     lipgloss.NewStyle().Background(lipgloss.Color("#444444")).Italic(true),
	}
}
