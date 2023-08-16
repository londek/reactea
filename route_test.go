package reactea

import (
	"reflect"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestRoutePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but it didn't")
		} else {
			if r != "tried updating global route not in update" {
				t.Errorf("expected panic, got it, but with invalid message, got \"%s\", expected \"tried updating global route not in update\"", r)
			}
		}
	}()

	NewProgram(&testComponenent{}, tea.WithoutRenderer())

	SetRoute("/shouldFail")
}

type testNavigateComponenent struct {
	BasicComponent
	BasicPropfulComponent[NoProps]

	routeHistory []string
	step         int
}

func (c *testNavigateComponenent) Init(props NoProps) tea.Cmd {
	c.UpdateProps(props)

	return Rerender
}

// Expecting:
// / -> /foo -> /foo/bar -> /baz -> / -> /foo -> /bar -> /test -> / -> / -> /foo -> /bar
func (c *testNavigateComponenent) Update(msg tea.Msg) tea.Cmd {
	switch c.step {
	case 0:
		// Don't navigate
	case 1:
		Navigate("foo")
	case 2:
		Navigate("foo/bar")
	case 3:
		Navigate("/baz")
	case 4:
		Navigate("..")
	case 5:
		Navigate("./foo")
	case 6:
		Navigate(".//bar")
	case 7:
		Navigate("../../test")
	case 8:
		Navigate("/")
	case 9:
		Navigate("")
	case 10:
		Navigate(".")
	case 11:
		Navigate("foo")
	case 12:
		Navigate("foo/bar")
	case 13:
		Navigate("baz")
	default:
		return Destroy
	}

	c.routeHistory = append(c.routeHistory, CurrentRoute())

	c.step += 1

	return Rerender
}

func (c *testNavigateComponenent) Render(width int, height int) string {
	return ""
}

func TestNavigate(t *testing.T) {
	root := &testNavigateComponenent{}
	program := NewProgram(root, tea.WithoutRenderer(), WithoutInput())

	if _, err := program.Run(); err != nil {
		t.Fatal(err)
	}

	expectedRouteHistory := []string{
		"/",
		"/foo",
		"/foo/bar",
		"/baz",
		"/",
		"/foo",
		"/bar",
		"/test",
		"/",
		"/",
		"/",
		"/foo",
		"/foo/bar",
		"/foo/baz",
	}

	if strings.Join(root.routeHistory, " - ") != strings.Join(expectedRouteHistory, " - ") {
		t.Errorf("wrong route history, expected \"%s\", got \"%s\". Note that routes are delimited by \" - \"", strings.Join(expectedRouteHistory, " - "), strings.Join(root.routeHistory, " - "))
	}
}

func TestRoutePlaceholderMatching(t *testing.T) {
	testCases := []struct {
		route, placeholder string
		expected           map[string]string
	}{
		// Matching against non-root routes is forbidden
		{"", "", nil},
		{"invalidRoute", "", nil},
		{"", "invalidPlaceholder", nil},
		{"/invalidRoute", "invalidPlaceholder", nil},

		{"/teams/foo", "/teams", nil},
		{"/teams", "/teams/foo", nil},
		{"/", "/teams", nil},
		{"/teams", "/", nil},
		{"/teams", "/teams", map[string]string{"$": "/teams"}},

		{"/teams", "/teams/?:", map[string]string{"$": "/teams"}},
		{"/teams/123", "/teams/?:", map[string]string{"$": "/teams/123"}},
		{"/teams/123/456", "/teams/?:", nil},
		{"/teams", "/teams/?:teamId", map[string]string{"$": "/teams", "teamId": ""}},
		{"/teams/123", "/teams/?:teamId", map[string]string{"$": "/teams/123", "teamId": "123"}},
		{"/teams/123/456", "/teams/?:teamId", nil},

		{"/teams/123/456", "/teams/123/456/+?:foo", map[string]string{"$": "/teams/123/456", "foo": ""}},
		{"/teams/123/456", "/teams/+?:foo", map[string]string{"$": "/teams/123/456", "foo": "123/456"}},
		{"/teams/123/456", "/teams/+?:", map[string]string{"$": "/teams/123/456"}},
		{"/teams/123", "/teams/+?:", map[string]string{"$": "/teams/123"}},
		{"/teams", "/teams/+?:", map[string]string{"$": "/teams"}},

		{"/teams/123", "/teams/:teamId", map[string]string{"$": "/teams/123", "teamId": "123"}},
		{"/teams/foo/234", "/teams/:/:teamId", map[string]string{"$": "/teams/foo/234", "teamId": "234"}},
		{"/teams/123/234", "/teams/:teamId/:teamId", map[string]string{"$": "/teams/123/234", "teamId": "234"}},
		{"/teams/123/234", "/teams/:teamId/:playerId", map[string]string{"$": "/teams/123/234", "teamId": "123", "playerId": "234"}},
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
