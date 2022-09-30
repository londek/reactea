package reactea

import "testing"

func TestOptions(t *testing.T) {
	t.Run("WithRoute", func(t *testing.T) {
		NewProgram(&testComponenent{}, WithRoute("testRoute"))

		if CurrentRoute() != "testRoute" {
			t.Errorf("expected current route \"testRoute\", but got \"%s\"", CurrentRoute())
		}
	})
}
