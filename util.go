package reactea

// There are no unions yet in Golang (Currently they are Type constraints BUT
// type constraints can't be used as argument type)
func RenderAny[TProps any](renderer any, props TProps, width, height int) string {
	if renderer, ok := renderer.(Renderer[TProps]); ok {
		return renderer(props, width, height)
	}

	return RenderPropless(renderer, width, height)
}

func RenderPropless(renderer any, width, height int) string {
	switch renderer := renderer.(type) {
	case ProplessRenderer:
		return renderer(width, height)
	case DumbRenderer:
		return renderer()
	}

	return ""
}

// Wrapper function
func PropfulToLess[TProps any](renderer Renderer[TProps], props TProps) ProplessRenderer {
	return func(width, height int) string {
		return renderer(props, width, height)
	}
}

// Transformer for Propless -> Some component

type proplessWrapper struct {
	BasicComponent

	Renderer ProplessRenderer
}

func (c *proplessWrapper) Render(width, height int) string {
	return c.Renderer(width, height)
}

func ProplessToSomeComponent(renderer ProplessRenderer) SomeComponent {
	return &proplessWrapper{Renderer: renderer}
}
