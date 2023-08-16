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

	routesHistory []string
	step          int
}

func (c *testNavigateComponenent) Init(props NoProps) tea.Cmd {
	c.UpdateProps(props)

	return Rerender
}

// Expecting:
// / -> /foo -> /foo/bar -> /baz -> / -> /foo -> /foo/bar -> /test -> /test -> /test -> /foo -> /bar
func (c *testNavigateComponenent) Update(msg tea.Msg) tea.Cmd {
	c.routesHistory = append(c.routesHistory, CurrentRoute())

	switch c.step {
	case 0:
		Navigate("foo")
	case 1:
		Navigate("foo/bar")
	case 2:
		Navigate("/baz")
	case 3:
		Navigate("..")
	case 4:
		Navigate("./foo")
	case 5:
		Navigate(".//bar")
	case 6:
		Navigate("../../test")
	case 7:
		Navigate("/")
	case 8:
		Navigate("")
	case 9:
		Navigate(".")
	case 10:
		Navigate("foo")
	case 11:
		Navigate("bar")
	default:
		return Destroy
	}

	c.step += 1

	return Rerender
}

func (c *testNavigateComponenent) Render(width int, height int) string {
	return ""
}

func TestNavigate(t *testing.T) {
	root := &testNavigateComponenent{}
	program := NewProgram(root, tea.WithoutRenderer(), WithoutInput())

	if err := program.Start(); err != nil {
		t.Fatal(err)
	}

	t.Log(strings.Join(root.routesHistory, " - "))
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
