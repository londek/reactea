package reactea

import tea "github.com/charmbracelet/bubbletea"

type BeforeUpdater interface {
	BeforeUpdate() tea.Cmd
}

var beforeUpdaters []BeforeUpdater

// Queue component for AfterUpdate event
func BeforeUpdate(beforeUpdater BeforeUpdater) {
	if beforeUpdater == nil {
		return
	}

	beforeUpdaters = append(beforeUpdaters, beforeUpdater)
}

func handleBeforeUpdates() tea.Cmd {
	if beforeUpdaters == nil {
		// Meaning it hasn't been updated, we could
		// len(afterUpdaters) == 0 but there is no
		// reason to because it will either be nil
		// or slice with elements
		return nil
	}

	cmds := make([]tea.Cmd, len(beforeUpdaters))

	for i, beforeUpdater := range beforeUpdaters {
		cmds[i] = beforeUpdater.BeforeUpdate()
	}

	beforeUpdaters = nil

	return tea.Batch(cmds...)
}
