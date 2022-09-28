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

	firstRun bool

	testInitializer func() map[string]RouteInitializer
	testUpdater     func(bool) tea.Cmd

	router *Component
}

func (c *testComponenent) Init(reactea.NoProps) tea.Cmd {
	return c.router.Init(c.testInitializer())
}

func (c *testComponenent) Update(msg tea.Msg) tea.Cmd {
	defer func() {
		c.firstRun = false
	}()

	cmd := c.testUpdater(c.firstRun)

	return tea.Batch(cmd, c.router.Update(msg))
}

func (c *testComponenent) Render(width, height int) string {
	return c.router.Render(width, height)
}

func TestDefault(t *testing.T) {
	var in, out bytes.Buffer

	in.WriteString("123")

	root := &testComponenent{
		firstRun: true,
		testInitializer: func() map[string]RouteInitializer {
			return map[string]RouteInitializer{
				"default": func() (reactea.SomeComponent, tea.Cmd) {
					renderer := func() string {
						return "Hello Default!"
					}

					return reactea.SomeComponentify(renderer, reactea.NoProps{}), nil
				},
			}
		},
		testUpdater: func(b bool) tea.Cmd {
			return reactea.Destroy
		},
		router: New(),
	}

	program := tea.NewProgram(reactea.New(root), tea.WithInput(&in), tea.WithOutput(&out))

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
		firstRun: true,
		testInitializer: func() map[string]RouteInitializer {
			return map[string]RouteInitializer{
				"default": func() (reactea.SomeComponent, tea.Cmd) {
					renderer := func() string {
						return "Hello Default!"
					}

					return reactea.SomeComponentify(renderer, reactea.NoProps{}), nil
				},
				"test/test": func() (reactea.SomeComponent, tea.Cmd) {
					renderer := func() string {
						return "Hello Tests!"
					}

					return reactea.SomeComponentify(renderer, reactea.NoProps{}), nil
				},
			}
		},
		testUpdater: func(b bool) tea.Cmd {
			if b {
				reactea.SetCurrentRoute("test/test")
				return nil
			} else {
				return reactea.Destroy
			}
		},
		router: New(),
	}

	program := tea.NewProgram(reactea.New(root), tea.WithInput(&in), tea.WithOutput(&out))

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
		firstRun: true,
		testInitializer: func() map[string]RouteInitializer {
			return map[string]RouteInitializer{}
		},
		testUpdater: func(b bool) tea.Cmd {
			return reactea.Destroy
		},
		router: New(),
	}

	program := tea.NewProgram(reactea.New(root), tea.WithInput(&in), tea.WithOutput(&out))

	if err := program.Start(); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out.String(), "Couldn't route for") {
		t.Fatalf("got invalid route message")
	}
}
