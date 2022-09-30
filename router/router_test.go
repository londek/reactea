package router

import (
	"bytes"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type testComponenent struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	router *Component

	testRoutes  map[string]RouteInitializer
	testUpdater func(*testComponenent) tea.Cmd

	updateN int
}

func (c *testComponenent) Init(reactea.NoProps) tea.Cmd {
	return c.router.Init(c.testRoutes)
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
		testRoutes: map[string]RouteInitializer{
			"default": func() (reactea.SomeComponent, tea.Cmd) {
				renderer := func() string {
					return "Hello Default!"
				}

				return reactea.Componentify[reactea.NoProps](renderer), nil
			},
		},
		router: New(),
	}

	program := reactea.NewProgram(root, tea.WithInput(&in), tea.WithOutput(&out))

	if err := program.Start(); err != nil {
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
		testRoutes: map[string]RouteInitializer{
			"default": func() (reactea.SomeComponent, tea.Cmd) {
				renderer := func() string {
					return "Hello Default!"
				}

				return reactea.Componentify[reactea.NoProps](renderer), nil
			},
			"test/test": func() (reactea.SomeComponent, tea.Cmd) {
				renderer := func() string {
					return "Hello Tests!"
				}

				return reactea.Componentify[reactea.NoProps](renderer), nil
			},
		},
		router: New(),
	}

	program := reactea.NewProgram(root, reactea.WithRoute("test/test"), tea.WithInput(&in), tea.WithOutput(&out))

	if err := program.Start(); err != nil {
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
		testRoutes: map[string]RouteInitializer{
			"default": func() (reactea.SomeComponent, tea.Cmd) {
				renderer := func() string {
					return "Hello Default!"
				}

				return reactea.Componentify[reactea.NoProps](renderer), nil
			},
			"test/test": func() (reactea.SomeComponent, tea.Cmd) {
				renderer := func() string {
					return "Hello Tests!"
				}

				return reactea.Componentify[reactea.NoProps](renderer), nil
			},
		},
		testUpdater: func(c *testComponenent) tea.Cmd {
			if c.updateN == 0 {
				reactea.SetCurrentRoute("test/test")

				return nil
			} else {
				return reactea.Destroy
			}
		},
		router: New(),
	}

	program := reactea.NewProgram(root, tea.WithInput(&in), tea.WithOutput(&out))

	if err := program.Start(); err != nil {
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

	if err := program.Start(); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out.String(), "Couldn't route for") {
		t.Fatalf("got invalid route message")
	}
}
