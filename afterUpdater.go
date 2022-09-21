package reactea

import tea "github.com/charmbracelet/bubbletea"

var afterUpdaters []AfterUpdater

// Queue component for AfterUpdate event
func AfterUpdate(afterUpdater AfterUpdater) {
	if afterUpdater == nil {
		return
	}

	afterUpdaters = append(afterUpdaters, afterUpdater)
}

func handleAfterUpdates() tea.Cmd {
	cmds := make([]tea.Cmd, len(afterUpdaters))

	for i, afterUpdater := range afterUpdaters {
		cmds[i] = afterUpdater.AfterUpdate()
	}

	afterUpdaters = nil

	return tea.Batch(cmds...)
}
