package reactea

import (
	"bytes"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type testBeforeUpdaterComponenent struct {
	BasicComponent

	firstRun bool
	text     string
}

func (c *testBeforeUpdaterComponenent) Update(msg tea.Msg) tea.Cmd {
	if c.firstRun {
		c.firstRun = false
		BeforeUpdate(c)
		return nil
	}

	return Destroy
}

func (c *testBeforeUpdaterComponenent) BeforeUpdate() tea.Cmd {
	c.text = "Hello Tests!"
	return nil
}

func (c *testBeforeUpdaterComponenent) Render(int, int) string {
	return c.text
}

func TestBeforeUpdate(t *testing.T) {
	var in, out bytes.Buffer

	in.WriteString("123")

	root := &testBeforeUpdaterComponenent{
		firstRun: true,
		text:     "Bad Test!",
	}

	program := NewProgram(root, tea.WithInput(&in), tea.WithOutput(&out))

	if _, err := program.Run(); err != nil {
		t.Fatal(err)
	}

	if strings.Contains(out.String(), "Bad Test!") {
		t.Errorf("got bad test")
	}

	if !strings.Contains(out.String(), "Hello Tests!") {
		t.Errorf("invalid render")
	}
}

func TestBeforeUpdateNil(t *testing.T) {
	BeforeUpdate(nil)

	if len(beforeUpdaters) != 0 {
		t.Errorf("expected beforeUpdaters slice to have len 0, but got %d", len(beforeUpdaters))
	}
}
