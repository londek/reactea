package reactea

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type testComponenent struct {
	BasicComponent
	BasicPropfulComponent[NoProps]

	echoKey string
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

func (c *testComponenent) Render(int, int) string {
	return c.echoKey
}

func TestComponent(t *testing.T) {
	var in, out bytes.Buffer

	in.WriteString("~~~x")

	root := &testComponenent{
		echoKey: "default",
	}

	program := tea.NewProgram(New(root), tea.WithInput(&in), tea.WithOutput(&out))

	if err := program.Start(); err != nil {
		t.Fatal(err)
	}

	if strings.Contains(out.String(), "default") {
		t.Errorf("TestComponent did not echo!")
	}

	if !strings.Contains(out.String(), "~") {
		t.Errorf("TestComponent invalid echo!")
	}

	if WasRouteChanged() {
		t.Errorf("TestComponent current route was changed!")
	}

	if CurrentRoute() != "test/test/test" {
		t.Errorf("TestComponent current route does not equal \"test/test/test\", got \"%s\"", CurrentRoute())
	}

	if props := root.Props(); !reflect.DeepEqual(props, NoProps{}) {
		t.Errorf("TestComponent props does not equal zero NoProps, got \"%s\"", props)
	}
}

func TestInvisibleComponent(t *testing.T) {
	component := &InvisibleComponent{}

	if result := component.Render(1, 1); result != "" {
		t.Errorf("TestInvisibleComponent expected empty string, got \"%s\"", result)
	}
}
