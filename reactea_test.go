package reactea

import (
	"bytes"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestComponent(t *testing.T) {
	var in, out bytes.Buffer

	in.WriteString("~~~")

	type testState struct {
		echoKey               string
		lastWidth, lastHeight int
	}

	root := &mockComponent[testState]{
		updateFunc: func(c Component, s *testState, msg tea.Msg) tea.Cmd {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				if msg.String() == "x" {
					return Destroy
				}

				s.echoKey = msg.String()
			}

			SetRoute("/test/test/test")

			return nil
		},
		renderFunc: func(c Component, s *testState, width, height int) string {
			s.lastWidth, s.lastHeight = width, height

			return s.echoKey
		},
	}

	program := NewProgram(root, tea.WithInput(&in), tea.WithOutput(&out))

	// Test for window size
	go func() {
		// Simulate initial window size
		program.Send(tea.WindowSizeMsg{Width: 1, Height: 1})

		// Give time to catch up
		time.Sleep(50 * time.Millisecond)

		// Simulate pressing X
		program.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}, Alt: false})
	}()

	if _, err := program.Run(); err != nil {
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

	if CurrentRoute() != "/test/test/test" {
		t.Errorf("current route is wrong, expected \"/test/test/test\", got \"%s\"", CurrentRoute())
	}

	if root.state.lastWidth != 1 {
		t.Errorf("expected lastWidth 1, but got %d", root.state.lastWidth)
	}

	if root.state.lastHeight != 1 {
		t.Errorf("expected lastHeigth 1, but got %d", root.state.lastWidth)
	}
}

func TestNew(t *testing.T) {
	t.Run("NewProgram", func(t *testing.T) {
		root := &mockComponent[struct{}]{
			renderFunc: func(c Component, s *struct{}, width, height int) string {
				return "test passed"
			},
		}

		program := NewProgram(root, WithoutInput(), tea.WithoutRenderer())

		go program.Quit()

		if _, err := program.Run(); err != nil {
			t.Fatal(err)
		}
	})
}
