package reactea

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type RouteUpdatedMsg struct {
	Original string
}

// Global Route object, substitute of window.location
// Feel free to use standard path package
var (
	currentRoute    = "/"
	lastRoute       = "/"
	wasRouteChanged = false
)

func CurrentRoute() string {
	return currentRoute
}

func LastRoute() string {
	return lastRoute
}

func WasRouteChanged() bool {
	return wasRouteChanged
}

func SetRoute(target string) {
	if !isUpdate {
		panic("tried updating global route not in update")
	}

	if len(target) == 0 || target[0] != '/' {
		panic("can't set route to non-root path")
	}

	if !wasRouteChanged {
		wasRouteChanged = currentRoute != target
		lastRoute = currentRoute
	}

	currentRoute = target
}

func Navigate(target string) {
	var currentRouteLevels []string

	if len(target) == 0 {
		// Just don't navigate if no target was given
		return
	}

	// Check whether target is absolute
	if target[0] == '/' {
		currentRouteLevels = []string{}
	} else {
		currentRouteLevels = strings.Split(currentRoute, "/")
		currentRouteLevels = currentRouteLevels[1 : len(currentRouteLevels)-1]
	}

	for _, targetLevel := range strings.Split(target, "/") {
		switch targetLevel {
		case ".", "":
		case "..":
			if len(currentRouteLevels) > 0 {
				currentRouteLevels = currentRouteLevels[:len(currentRouteLevels)-1]
			}
		default:
			currentRouteLevels = append(currentRouteLevels, targetLevel)
		}
	}

	SetRoute("/" + strings.Join(currentRouteLevels, "/"))
}

// Checks whether route (e.g. /teams/123/12) matches
// placeholder (e.g /teams/:teamId/:playerId)
// Returns map of found params and if it matches
// Params have to follow regex ^:.*$
// ^ being beginning of current path level (/^level/)
// $ being end of current path level (/level$/)
//
// Note: Entire matched route can be accessed with key "$"
// Note: Placeholders can be optional => /foo/?:/?: will match foo/bar and foo and /foo/bar/baz
// Note: The most outside placeholders can be optional recursive => /foo/+?: will match /foo/bar and foo and /foo/bar/baz
// Note: It allows for defining wildcards with /foo/:/bar
// Note: Duplicate params will result in overwrite of first param
func RouteMatchesPlaceholder(route string, placeholder string) (map[string]string, bool) {
	var (
		routeLevels       = strings.Split(route, "/")
		placeholderLevels = strings.Split(placeholder, "/")
	)

	if len(route) > 0 && route[0] == '/' {
		routeLevels = routeLevels[1:]
	} else {
		// Checking against non-root route is forbidden
		return nil, false
	}

	if len(placeholder) > 0 && placeholder[0] == '/' {
		placeholderLevels = placeholderLevels[1:]
	} else {
		// Checking against non-root placeholder is forbidden
		return nil, false
	}

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
func updatedRoute(original string) tea.Cmd {
	return func() tea.Msg {
		return RouteUpdatedMsg{original}
	}
}
