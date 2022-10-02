package reactea

import (
	"path"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

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

// Checks whether route (e.g. teams/123/12) matches
// placeholder (e.g teams/:teamId/:playerId)
// Returns map of found params and if it matches
// Params have to follow regex ^:.*$
// ^ being beginning of current path level (/^level/)
// $ being end of current path level (/level$/)
//
// Note: Entire matched route can be accessed with key "$"
// Note: that it allows for defining wildcards with foo/:/bar
// Note: Duplicate params will result in overwrite of first param
func RouteMatchesPlaceholder(route string, placeholder string) (params map[string]string, ok bool) {
	var (
		routeLevels       = strings.Split(path.Clean(route), "/")
		placeholderLevels = strings.Split(path.Clean(placeholder), "/")
	)

	if len(routeLevels) != len(placeholderLevels) {
		return
	}

	params = make(map[string]string, len(placeholderLevels)+1)

	params["$"] = route

	for i, routeLevel := range routeLevels {
		placeholderLevel := placeholderLevels[i]

		if len(placeholderLevel) > 0 && placeholderLevel[0] == ':' {
			if placeholderLevel == ":" {
				continue // wildcard
			}

			paramName := placeholderLevel[1:]
			params[paramName] = routeLevel
		} else {
			if routeLevel != placeholderLevel {
				return
			}
		}
	}

	ok = true

	return
}

// It might be important to do so in some scenarios.
// Basically causes rerender so ALL components are
// aware of changed routes
func updatedRoute() tea.Msg {
	return nil
}
