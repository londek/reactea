package reactea

import (
	"testing"
)

func TestRenderAny(t *testing.T) {
	t.Run("renderer", func(t *testing.T) {
		renderer := func(struct{}, int, int) string {
			return "working"
		}

		if result := RenderAny(renderer, struct{}{}, 1, 1); result != "working" {
			t.Errorf("invalid result, expected \"working\", got \"%s\"", result)
		}
	})

	t.Run("proplessRenderer", func(t *testing.T) {
		proplessRenderer := func(int, int) string {
			return "working"
		}

		if result := RenderAny(proplessRenderer, struct{}{}, 1, 1); result != "working" {
			t.Errorf("invalid result, expected \"working\", got \"%s\"", result)
		}
	})

	t.Run("dumbRenderer", func(t *testing.T) {
		dumbRenderer := func() string {
			return "working"
		}

		if result := RenderAny(dumbRenderer, struct{}{}, 1, 1); result != "working" {
			t.Errorf("invalid result, expected \"working\", got \"%s\"", result)
		}
	})
}

func TestPropfulToLess(t *testing.T) {
	renderer := func(struct{}, int, int) string {
		return "working"
	}

	proplessRenderer := PropfulToLess(renderer, struct{}{})

	if result := proplessRenderer(1, 1); result != "working" {
		t.Errorf("wrapped value doesn't render correctly, expected \"working\", got \"%s\"", result)
	}
}

func TestComponentify(t *testing.T) {
	t.Run("renderer", func(t *testing.T) {
		renderer := func(struct{}, int, int) string {
			return "working"
		}

		if result := Componentify[struct{}](renderer, struct{}{}).Render(1, 1); result != "working" {
			t.Errorf("transformed value doesn't render correctly, expected \"working\", got \"%s\"", result)
		}
	})

	t.Run("proplessRenderer", func(t *testing.T) {
		proplessRenderer := func(int, int) string {
			return "working"
		}

		if result := Componentify[struct{}](proplessRenderer, struct{}{}).Render(1, 1); result != "working" {
			t.Errorf("transformed value doesn't render correctly, expected \"working\", got \"%s\"", result)
		}
	})

	t.Run("dumbRenderer", func(t *testing.T) {
		dumbRenderer := func() string {
			return "working"
		}

		if result := Componentify[struct{}](dumbRenderer, struct{}{}).Render(1, 1); result != "working" {
			t.Errorf("transformed value doesn't render correctly, expected \"working\", got \"%s\"", result)
		}
	})
}
