package router

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type testComponenent struct {
	reactea.BasicComponent

	router *Component

	testUpdater func(*testComponenent) tea.Cmd

	updateN int
}

func (c *testComponenent) Init() tea.Cmd {
	return c.router.Init()
}

func (c *testComponenent) Update(msg tea.Msg) tea.Cmd {
	defer func() {
		c.updateN++
	}()

	if c.testUpdater != nil {
		return tea.Batch(c.router.Update(c.router.Update(msg)), c.testUpdater(c))
	}

	return tea.Batch(c.router.Update(msg), reactea.Destroy)
}

func (c *testComponenent) Render(width, height int) string {
	return c.router.Render(width, height)
}

func TestDefault(t *testing.T) {
	var in, out bytes.Buffer

	in.WriteString("123")

	root := &testComponenent{
		router: NewWithRoutes(map[string]RouteInitializer{
			"default": func(Params) reactea.Component {
				renderer := func() string {
					return "Hello Default!"
				}

				return reactea.ComponentifyDumb(renderer)
			},
		}),
	}

	program := reactea.NewProgram(root, tea.WithInput(&in), tea.WithOutput(&out))

	if _, err := program.Run(); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out.String(), "Hello Default!") {
		t.Fatalf("no default route message")
	}
}

func TestNonDefault(t *testing.T) {
	var in, out bytes.Buffer

	in.WriteString("123")

	root := &testComponenent{
		router: NewWithRoutes(map[string]RouteInitializer{
			"default": func(Params) reactea.Component {
				renderer := func() string {
					return "Hello Default!"
				}

				return reactea.ComponentifyDumb(renderer)
			},
			"/test/test": func(Params) reactea.Component {
				renderer := func() string {
					return "Hello Tests!"
				}

				return reactea.ComponentifyDumb(renderer)
			},
		}),
	}

	program := reactea.NewProgram(root, reactea.WithRoute("/test/test"), tea.WithInput(&in), tea.WithOutput(&out))

	if _, err := program.Run(); err != nil {
		t.Fatal(err)
	}

	if strings.Contains(out.String(), "Hello Default!") {
		t.Fatalf("got default route message")
	}

	if !strings.Contains(out.String(), "Hello Tests!") {
		t.Fatalf("got invalid route message")
	}
}

func TestRouteChange(t *testing.T) {
	var in, out bytes.Buffer

	in.WriteString("123")

	root := &testComponenent{
		testUpdater: func(c *testComponenent) tea.Cmd {
			if c.updateN == 0 {
				reactea.SetRoute("/test/test")

				return nil
			} else {
				return reactea.Destroy
			}
		},
		router: NewWithRoutes(map[string]RouteInitializer{
			"default": func(Params) reactea.Component {
				renderer := func() string {
					return "Hello Default!"
				}

				return reactea.ComponentifyDumb(renderer)
			},
			"/test/test": func(Params) reactea.Component {
				renderer := func() string {
					return "Hello Tests!"
				}

				return reactea.ComponentifyDumb(renderer)
			},
		}),
	}

	program := reactea.NewProgram(root, tea.WithInput(&in), tea.WithOutput(&out))

	if _, err := program.Run(); err != nil {
		t.Fatal(err)
	}

	if strings.Contains(out.String(), "Hello Default!") {
		t.Fatalf("got default route message")
	}

	if !strings.Contains(out.String(), "Hello Tests!") {
		t.Fatalf("got invalid route message")
	}
}

func TestNotFound(t *testing.T) {
	var in, out bytes.Buffer

	in.WriteString("123")

	root := &testComponenent{
		router: New(),
	}

	program := reactea.NewProgram(root, tea.WithInput(&in), tea.WithOutput(&out))

	if _, err := program.Run(); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out.String(), "Couldn't route for") {
		t.Fatalf("got invalid route message")
	}
}

func TestRouteWithParam(t *testing.T) {
	var in, out bytes.Buffer

	in.WriteString("123")

	root := &testComponenent{
		testUpdater: func(c *testComponenent) tea.Cmd {
			if c.updateN == 0 {
				reactea.SetRoute("/test/wellDone")

				return nil
			} else {
				return reactea.Destroy
			}
		},
		router: NewWithRoutes(map[string]RouteInitializer{
			"default": func(Params) reactea.Component {
				renderer := func() string {
					return "Hello Default!"
				}

				return reactea.ComponentifyDumb(renderer)
			},
			"/test/:foo": func(params Params) reactea.Component {
				renderer := func() string {
					return fmt.Sprintf("Hello Tests! Param foo is %s", params["foo"])
				}

				return reactea.ComponentifyDumb(renderer)
			},
		}),
	}

	program := reactea.NewProgram(root, tea.WithInput(&in), tea.WithOutput(&out))

	if _, err := program.Run(); err != nil {
		t.Fatal(err)
	}

	if strings.Contains(out.String(), "Hello Default!") {
		t.Fatalf("got default route message")
	}

	if !strings.Contains(out.String(), "Hello Tests!") {
		t.Fatalf("got invalid route message")
	}

	if !strings.Contains(out.String(), "Hello Tests! Param foo is wellDone") {
		t.Fatalf("got valid route message, but most likely wrong param")
	}
}
