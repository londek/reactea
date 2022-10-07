package reactea

import (
	"reflect"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestRoutePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but it didn't")
		}
	}()

	NewProgram(&testComponenent{}, tea.WithoutRenderer())

	SetCurrentRoute("shouldFail")
}

func TestRoutePlaceholderMatching(t *testing.T) {
	testCases := []struct {
		route, placeholder string
		expected           map[string]string
	}{
		{"", "", map[string]string{"$": ""}},
		{"teams/foo", "teams", nil},
		{"", "teams", nil},
		{"teams", "", nil},
		{"teams", "teams", map[string]string{"$": "teams"}},

		{"teams", "teams/?:", map[string]string{"$": "teams"}},
		{"teams/123", "teams/?:", map[string]string{"$": "teams/123"}},
		{"teams/123/456", "teams/?:", nil},
		{"teams", "teams/?:teamId", map[string]string{"$": "teams", "teamId": ""}},
		{"teams/123", "teams/?:teamId", map[string]string{"$": "teams/123", "teamId": "123"}},
		{"teams/123/456", "teams/?:teamId", nil},

		{"teams/123/456", "teams/123/456/+?:foo", map[string]string{"$": "teams/123/456", "foo": ""}},
		{"teams/123/456", "teams/+?:foo", map[string]string{"$": "teams/123/456", "foo": "123/456"}},
		{"teams/123/456", "teams/+?:", map[string]string{"$": "teams/123/456"}},
		{"teams/123", "teams/+?:", map[string]string{"$": "teams/123"}},
		{"teams", "teams/+?:", map[string]string{"$": "teams"}},

		{"teams/123", "teams/:teamId", map[string]string{"$": "teams/123", "teamId": "123"}},
		{"teams/foo/234", "teams/:/:teamId", map[string]string{"$": "teams/foo/234", "teamId": "234"}},
		{"teams/123/234", "teams/:teamId/:teamId", map[string]string{"$": "teams/123/234", "teamId": "234"}},
		{"teams/123/234", "teams/:teamId/:playerId", map[string]string{"$": "teams/123/234", "teamId": "123", "playerId": "234"}},
	}

	for _, testCase := range testCases {
		got, ok := RouteMatchesPlaceholder(testCase.route, testCase.placeholder)

		if testCase.expected == nil {
			if !ok {
				continue
			}

			t.Errorf("Bad result. Route: \"%s\", Placeholder: \"%s\". Expected not ok, got ok", testCase.route, testCase.placeholder)
			continue
		}

		if !reflect.DeepEqual(got, testCase.expected) {
			t.Errorf("Bad result. Route: \"%s\", Placeholder: \"%s\". Expected %v, got %v", testCase.route, testCase.placeholder, testCase.expected, got)
		}
	}
}
