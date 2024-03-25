package app

func (m Model) View() (out string) {
	if m.quitting {
		return ""
	}
	return m.Context + " --- " + m.Filters.Read
}
