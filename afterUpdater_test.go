package reactea

import (
	"bytes"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type testAfterUpdaterComponenent struct {
	BasicComponent
	BasicPropfulComponent[NoProps]

	firstRun bool
	text     string
}

func (c *testAfterUpdaterComponenent) Update(msg tea.Msg) tea.Cmd {
	if c.firstRun {
		c.firstRun = false
		AfterUpdate(c)
		return nil
	}

	return Destroy
}

func (c *testAfterUpdaterComponenent) AfterUpdate() tea.Cmd {
	c.text = "Hello Tests!"

	return nil
}

func (c *testAfterUpdaterComponenent) Render(int, int) string {
	return c.text
}

func TestAfterUpdater(t *testing.T) {
	var in, out bytes.Buffer

	in.WriteString("123")

	root := &testAfterUpdaterComponenent{
		firstRun: true,
		text:     "Bad Test!",
	}

	program := tea.NewProgram(New(root), tea.WithInput(&in), tea.WithOutput(&out))

	if err := program.Start(); err != nil {
		t.Fatal(err)
	}

	if strings.Contains(out.String(), "Bad Test!") {
		t.Errorf("TestAfterUpdater got bad test")
	}

	if !strings.Contains(out.String(), "Hello Tests!") {
		t.Errorf("TestAfterUpdater invalid render")
	}
}
