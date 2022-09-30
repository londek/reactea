package reactea

import tea "github.com/charmbracelet/bubbletea"

func WithRoute(route string) func(*tea.Program) {
	return func(*tea.Program) {
		currentRoute = route
	}
}
