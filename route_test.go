package reactea

import (
	"reflect"
	"testing"
)

func TestRoutePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but it didn't")
		}
	}()

	NewProgram(&testComponenent{})

	SetCurrentRoute("shouldFail")
}

func TestRoutePlaceholderMatching(t *testing.T) {
	testCases := []struct {
		route, placeholder string
		expected           map[string]string
	}{
		{"", "", map[string]string{}},
		{"teams/foo", "teams", nil},
		{"", "teams", nil},
		{"teams", "", nil},
		{"teams", "teams", map[string]string{}},
		{"teams/123", "teams/:teamId", map[string]string{"teamId": "123"}},
		{"teams/foo/234", "teams/:/:teamId", map[string]string{"teamId": "234"}},
		{"teams/123/234", "teams/:teamId/:teamId", map[string]string{"teamId": "234"}},
		{"teams/123/234", "teams/:teamId/:playerId", map[string]string{"teamId": "123", "playerId": "234"}},
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
