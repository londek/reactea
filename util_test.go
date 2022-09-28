package reactea

import (
	"testing"
)

func TestRenderAny(t *testing.T) {
	t.Run("renderer", func(t *testing.T) {
		renderer := func(NoProps, int, int) string {
			return "working"
		}

		if result := RenderAny(renderer, NoProps{}, 1, 1); result != "working" {
			t.Errorf("invalid result, expected \"working\", got \"%s\"", result)
		}
	})

	t.Run("proplessRenderer", func(t *testing.T) {
		proplessRenderer := func(int, int) string {
			return "working"
		}

		if result := RenderAny(proplessRenderer, NoProps{}, 1, 1); result != "working" {
			t.Errorf("invalid result, expected \"working\", got \"%s\"", result)
		}
	})

	t.Run("dumbRenderer", func(t *testing.T) {
		dumbRenderer := func() string {
			return "working"
		}

		if result := RenderAny(dumbRenderer, NoProps{}, 1, 1); result != "working" {
			t.Errorf("invalid result, expected \"working\", got \"%s\"", result)
		}
	})
}

func TestPropfulToLess(t *testing.T) {
	renderer := func(NoProps, int, int) string {
		return "working"
	}

	proplessRenderer := PropfulToLess(renderer, NoProps{})

	if result := proplessRenderer(1, 1); result != "working" {
		t.Errorf("wrapped value doesn't render correctly, expected \"working\", got \"%s\"", result)
	}
}

func TestComponentify(t *testing.T) {
	t.Run("renderer", func(t *testing.T) {
		renderer := func(NoProps, int, int) string {
			return "working"
		}

		if result := Componentify[NoProps](renderer).Render(1, 1); result != "working" {
			t.Errorf("transformed value doesn't render correctly, expected \"working\", got \"%s\"", result)
		}
	})

	t.Run("proplessRenderer", func(t *testing.T) {
		proplessRenderer := func(int, int) string {
			return "working"
		}

		if result := Componentify[NoProps](proplessRenderer).Render(1, 1); result != "working" {
			t.Errorf("transformed value doesn't render correctly, expected \"working\", got \"%s\"", result)
		}
	})

	t.Run("dumbRenderer", func(t *testing.T) {
		dumbRenderer := func() string {
			return "working"
		}

		if result := Componentify[NoProps](dumbRenderer).Render(1, 1); result != "working" {
			t.Errorf("transformed value doesn't render correctly, expected \"working\", got \"%s\"", result)
		}
	})
}

func TestSomeComponentify(t *testing.T) {
	t.Run("renderer", func(t *testing.T) {
		renderer := func(NoProps, int, int) string {
			return "working"
		}

		if result := SomeComponentify(renderer, NoProps{}).Render(1, 1); result != "working" {
			t.Errorf("transformed value doesn't render correctly, expected \"working\", got \"%s\"", result)
		}
	})

	t.Run("proplessRenderer", func(t *testing.T) {
		proplessRenderer := func(int, int) string {
			return "working"
		}

		if result := SomeComponentify(proplessRenderer, NoProps{}).Render(1, 1); result != "working" {
			t.Errorf("transformed value doesn't render correctly, expected \"working\", got \"%s\"", result)
		}
	})

	t.Run("dumbRenderer", func(t *testing.T) {
		dumbRenderer := func() string {
			return "working"
		}

		if result := SomeComponentify(dumbRenderer, NoProps{}).Render(1, 1); result != "working" {
			t.Errorf("transformed value doesn't render correctly, expected \"working\", got \"%s\"", result)
		}
	})
}
