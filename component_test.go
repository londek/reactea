package reactea

import (
	"bytes"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Component fully implemented by embedding premade structs
type testDefaultComponent struct {
	BasicComponent[NoProps]
}

func (*testDefaultComponent) Render(int, int) string {
	return "test passed"
}
func TestDefaultComponent(t *testing.T) {
	var out bytes.Buffer

	component := &testDefaultComponent{}

	program := NewProgram(component, WithoutInput(), tea.WithOutput(&out))

	go func() {
		t.Run("afterUpdate", func(t *testing.T) {
			AfterUpdate(component)

			// Force running update
			program.Send(nil)
		})

		time.Sleep(20 * time.Millisecond)

		program.Quit()
	}()

	if _, err := program.Run(); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out.String(), "test passed") {
		t.Errorf("invalid output, got \"%s\"", out.String())
	}
}

func TestInvisibleComponent(t *testing.T) {
	component := &InvisibleComponent{}

	if result := component.Render(1, 1); result != "" {
		t.Errorf("expected empty string, got \"%s\"", result)
	}
}
