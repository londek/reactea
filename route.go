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
// Note: Placeholders can be optional => foo/?:/?: will match foo/bar and foo and foo/bar/baz
// Note: The most outside placeholders can be optional recursive => foo/+?: will match foo/bar and foo and foo/bar/baz
// Note: It allows for defining wildcards with foo/:/bar
// Note: Duplicate params will result in overwrite of first param
func RouteMatchesPlaceholder(route string, placeholder string) (map[string]string, bool) {
	var (
		routeLevels       = strings.Split(path.Clean(route), "/")
		placeholderLevels = strings.Split(path.Clean(placeholder), "/")
	)

	if len(routeLevels) > len(placeholderLevels) && !strings.HasPrefix(placeholderLevels[len(placeholderLevels)-1], "+?:") {
		return nil, false
	}

	params := make(map[string]string, len(placeholderLevels)+1)

	params["$"] = route

	for i, placeholderLevel := range placeholderLevels {
		if i > len(routeLevels)-1 {
			if strings.HasPrefix(placeholderLevel, "?:") {
				if placeholderLevel == "?:" {
					continue // wildcard
				}

				paramName := placeholderLevel[2:]
				params[paramName] = ""
				continue
			} else if strings.HasPrefix(placeholderLevel, "+?:") {
				if placeholderLevel == "+?:" {
					break
				}

				paramName := placeholderLevel[3:]
				params[paramName] = ""
				break
			} else {
				// We are out of bounds and placeholder doesn't want optional data
				return nil, false
			}
		}

		routeLevel := routeLevels[i]

		if strings.HasPrefix(placeholderLevel, ":") {
			if placeholderLevel == ":" {
				continue // wildcard
			}

			paramName := placeholderLevel[1:]
			params[paramName] = routeLevel
		} else if strings.HasPrefix(placeholderLevel, "?:") {
			if placeholderLevel == "?:" {
				continue // wildcard
			}

			paramName := placeholderLevel[2:]
			params[paramName] = routeLevel
		} else if strings.HasPrefix(placeholderLevel, "+?:") {
			if placeholderLevel == "+?:" {
				break
			}

			paramName := placeholderLevel[3:]
			params[paramName] = strings.Join(routeLevels[i:], "/")
			break
		} else {
			if routeLevel != placeholderLevel {
				return nil, false
			}
		}
	}

	return params, true
}

// It might be important to do so in some scenarios.
// Basically causes rerender so ALL components are
// aware of changed routes
func updatedRoute() tea.Msg {
	return nil
}
