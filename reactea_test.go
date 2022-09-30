package reactea

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type testComponenent struct {
	BasicComponent
	BasicPropfulComponent[NoProps]

	echoKey string

	lastWidth, lastHeight int
}

func (c *testComponenent) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "x" {
			return Destroy
		}

		c.echoKey = msg.String()
	}

	SetCurrentRoute("test/test/test")

	return nil
}

func (c *testComponenent) Render(width int, height int) string {
	c.lastWidth, c.lastHeight = width, height

	return c.echoKey
}

func TestComponent(t *testing.T) {
	var in, out = &bytes.Buffer{}, &bytes.Buffer{}

	in.WriteString("~~~")

	root := &testComponenent{
		echoKey: "default",
	}

	program := NewProgram(root, tea.WithInput(in), tea.WithOutput(out))

	// Test for window size
	go func() {
		// Simulate initial window size
		program.Send(tea.WindowSizeMsg{Width: 1, Height: 1})

		// Give time to catch up
		time.Sleep(50 * time.Millisecond)

		// Simulate pressing X
		program.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}, Alt: false})
	}()

	if err := program.Start(); err != nil {
		t.Fatal(err)
	}

	if strings.Contains(out.String(), "default") {
		t.Errorf("did not echo")
	}

	if !strings.Contains(out.String(), "~") {
		t.Errorf("invalid echo")
	}

	if WasRouteChanged() {
		t.Errorf("current route was changed")
	}

	if CurrentRoute() != "test/test/test" {
		t.Errorf("current route is wrong, expected \"test/test/test\", got \"%s\"", CurrentRoute())
	}

	if props := root.Props(); !reflect.DeepEqual(props, NoProps{}) {
		t.Errorf("props is not zero-value of NoProps, got \"%s\"", props)
	}

	if root.lastWidth != 1 {
		t.Errorf("expected lastWidth 1, but got %d", root.lastWidth)
	}

	if root.lastHeight != 1 {
		t.Errorf("expected lastHeigth 1, but got %d", root.lastWidth)
	}
}

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		program := tea.NewProgram(New(&testDefaultComponent{}), WithoutInput(), tea.WithoutRenderer())

		go program.Quit()

		if err := program.Start(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("NewProgram", func(t *testing.T) {
		program := NewProgram(&testDefaultComponent{}, WithoutInput(), tea.WithoutRenderer())

		go program.Quit()

		if err := program.Start(); err != nil {
			t.Fatal(err)
		}
	})
}
