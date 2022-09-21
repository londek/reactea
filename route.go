package reactea

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Global Route object, substitute of window.location
var currentRoute Route

func CurrentRoute() Route {
	return currentRoute.Copy()
}

func SetCurrentRoute(r Route) {
	if !isUpdate {
		panic("tried updating global route in not update")
	}

	currentRoute = r.Copy()
}

type Route []string

func (r Route) Shift() (Route, Route) {
	if len(r) == 0 {
		return r, r
	}

	return r[1:], r[:1]
}

func (r Route) Pop() (Route, Route) {
	if len(r) == 0 {
		return r, r
	}

	return r[:len(r)-1], r[len(r)-1:]
}

func (r Route) Push(element string) Route {
	return append(r, element)
}

func (r Route) Copy() (dst Route) {
	return append(dst, r...)
}

func (r Route) String() string {
	return strings.Join(r, "/")
}

func (r1 Route) Equal(r2 Route) bool {
	if len(r1) != len(r2) {
		return false
	}

	for i, element := range r1 {
		if element != r2[i] {
			return false
		}
	}

	return true
}

func RouteOf(route string) (dst Route) {
	return strings.Split(route, "/")
}

type updatedRoutesMsg struct{}

// It might be important to do so in some scenarios.
// Basically causes rerender so ALL components are
// aware of changed routes
func updatedRoutes() tea.Msg {
	return updatedRoutesMsg{}
}
