package reactea

import tea "github.com/charmbracelet/bubbletea"

var afterUpdaters []Component

// Queue component for AfterUpdate event
func AfterUpdate(afterUpdater Component) {
	if afterUpdater == nil {
		return
	}

	afterUpdaters = append(afterUpdaters, afterUpdater)
}

func handleAfterUpdates() tea.Cmd {
	if afterUpdaters == nil {
		// Meaning it hasn't been updated, we could
		// len(afterUpdaters) == 0 but there is no
		// reason to because it will either be nil
		// or slice with elements
		return nil
	}

	cmds := make([]tea.Cmd, len(afterUpdaters))

	for i, afterUpdater := range afterUpdaters {
		cmds[i] = afterUpdater.AfterUpdate()
	}

	afterUpdaters = nil

	return tea.Batch(cmds...)
}
