package layout

// func cancelEdit(m Model) tea.Cmd {
// 	return func() tea.Msg {
// 		task := m.CurrentItem()
// 		m.input.Placeholder = task.Val()
// 		m.input.SetValue(task.Val())
// 		return EditCancelledMsg(m.input.Value()) //
// 	}
// }

// func (m Model) handleEdit(msg tea.Msg) (Model, tea.Cmd) {
// 	var cmd tea.Cmd
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.Type {
// 		case tea.KeyEnter:
// 			m.State = StateHome
// 			return m, changedEdit(m)
// 		case tea.KeyCtrlC, tea.KeyEsc:
// 			m.State = StateHome
// 			return m, cancelEdit(m)
// 		}
// 	case errMsg:
// 		m.err = msg
// 		return m, nil
// 	}
// 	m.input, cmd = m.input.Update(msg)
// 	return m, cmd
// }
