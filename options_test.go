package reactea

import "testing"

func TestOptions(t *testing.T) {
	t.Run("WithRoute", func(t *testing.T) {
		root := &mockComponent[struct{}]{}

		NewProgram(root, WithRoute("/testRoute"))

		if CurrentRoute() != "/testRoute" {
			t.Errorf("expected current route \"/testRoute\", but got \"%s\"", CurrentRoute())
		}
	})

	t.Run("WithoutInput", func(t *testing.T) {
		root := &mockComponent[struct{}]{
			renderFunc: func(c Component, s *struct{}, width, height int) string {
				return "test passed"
			},
		}

		program := NewProgram(root, WithoutInput())

		go program.Quit()

		if _, err := program.Run(); err != nil {
			t.Fatal(err)
		}
	})
}
