package reactea

import (
	"testing"
)

func TestRenderAny(t *testing.T) {
	renderer := func(NoProps, int, int) string {
		return "working"
	}

	proplessRenderer := func(int, int) string {
		return "working"
	}

	dumbRenderer := func() string {
		return "working"
	}

	if result := RenderAny(renderer, NoProps{}, 1, 1); result != "working" {
		t.Errorf("RenderAny failed for Renderer, got \"%s\"", result)
	}

	if result := RenderAny(proplessRenderer, NoProps{}, 1, 1); result != "working" {
		t.Errorf("RenderAny failed for ProplessRenderer, got \"%s\"", result)
	}

	if result := RenderAny(dumbRenderer, NoProps{}, 1, 1); result != "working" {
		t.Errorf("RenderAny failed for DumbRenderer, got \"%s\"", result)
	}
}

func TestPropfulToLess(t *testing.T) {
	renderer := func(NoProps, int, int) string {
		return "working"
	}

	proplessRenderer := PropfulToLess(renderer, NoProps{})

	if result := proplessRenderer(1, 1); result != "working" {
		t.Errorf("PropfulToLess wrapped value doesn't render correctly, got \"%s\"", result)
	}
}

func TestComponentify(t *testing.T) {
	renderer := func(NoProps, int, int) string {
		return "working"
	}

	proplessRenderer := func(int, int) string {
		return "working"
	}

	dumbRenderer := func() string {
		return "working"
	}

	if result := Componentify[NoProps](renderer).Render(1, 1); result != "working" {
		t.Errorf("Renderer Componentify transformed value doesn't render correctly, got \"%s\"", result)
	}

	if result := Componentify[NoProps](proplessRenderer).Render(1, 1); result != "working" {
		t.Errorf("ProplessRenderer Componentify transformed value doesn't render correctly, got \"%s\"", result)
	}

	if result := Componentify[NoProps](dumbRenderer).Render(1, 1); result != "working" {
		t.Errorf("DumbRenderer Componentify transformed value doesn't render correctly, got \"%s\"", result)
	}
}

func TestSomeComponentify(t *testing.T) {
	renderer := func(NoProps, int, int) string {
		return "working"
	}

	proplessRenderer := func(int, int) string {
		return "working"
	}

	dumbRenderer := func() string {
		return "working"
	}

	if result := SomeComponentify(renderer, NoProps{}).Render(1, 1); result != "working" {
		t.Errorf("Renderer Componentify transformed value doesn't render correctly, got \"%s\"", result)
	}

	if result := SomeComponentify(proplessRenderer, NoProps{}).Render(1, 1); result != "working" {
		t.Errorf("ProplessRenderer Componentify transformed value doesn't render correctly, got \"%s\"", result)
	}

	if result := SomeComponentify(dumbRenderer, NoProps{}).Render(1, 1); result != "working" {
		t.Errorf("DumbRenderer Componentify transformed value doesn't render correctly, got \"%s\"", result)
	}
}
