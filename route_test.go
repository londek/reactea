package reactea

import "testing"

func TestRoutePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but it didn't")
		}
	}()

	NewProgram(&testComponenent{})

	SetCurrentRoute("shouldFail")
}
