package reactea

var beforeUpdaters []Component

// Queue component for AfterUpdate event
func BeforeUpdate(beforeUpdater Component) {
	if beforeUpdater == nil {
		return
	}

	beforeUpdaters = append(beforeUpdaters, beforeUpdater)
}

func handleBeforeUpdates() {
	for _, beforeUpdater := range beforeUpdaters {
		beforeUpdater.BeforeUpdate()
	}

	beforeUpdaters = nil
}
