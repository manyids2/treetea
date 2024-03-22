package app

import (
	lgs "github.com/charmbracelet/lipgloss"
	ccc "github.com/manyids2/tasktea/components/theme"
)

type Styles struct {
	Normal lgs.Style
}

func NewStyles() Styles {
	return Styles{
		Normal: lgs.NewStyle().Foreground(lgs.Color(ccc.ColorMutedForeground)),
	}
}
