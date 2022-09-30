package reactea

import tea "github.com/charmbracelet/bubbletea"

// Global Route object, substitute of window.location
// Feel free to use standard path package
var (
	currentRoute    string
	wasRouteChanged bool
)

func CurrentRoute() string {
	return currentRoute
}

func WasRouteChanged() bool {
	return wasRouteChanged
}

func SetCurrentRoute(newRoute string) {
	if !isUpdate {
		panic("tried updating global route not in update")
	}

	if !wasRouteChanged {
		wasRouteChanged = currentRoute != newRoute
	}

	currentRoute = newRoute
}

// It might be important to do so in some scenarios.
// Basically causes rerender so ALL components are
// aware of changed routes
func updatedRoute() tea.Msg {
	return nil
}
