package colors

import lg "github.com/charmbracelet/lipgloss"

type Colors struct {
	Foreground lg.Color
	Background lg.Color
}

func New() Colors {
	return Colors{
		Foreground: lg.Color("#444444"),
		Background: lg.Color("#dddddd"),
	}
}
