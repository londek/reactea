package reactea

func RenderAny[TProps any, TRenderer AnyComponent[TProps]](renderer TRenderer, props TProps, width, height int) string {
	switch renderer := any(renderer).(type) {
	case Renderer[TProps]:
		return renderer(props, width, height)
	case ProplessRenderer:
		return renderer(width, height)
	case DumbRenderer:
		return renderer()
	}

	return ""
}

func RenderPropless[TRenderer AnyProplessComponent](renderer TRenderer, width, height int) string {
	switch renderer := any(renderer).(type) {
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

// Transformer for Propless -> SomeComponent

type proplessWrapper struct {
	BasicComponent

	Renderer ProplessRenderer
}

func (c *proplessWrapper) Render(width, height int) string {
	return c.Renderer(width, height)
}

func ProplessToComponent(renderer ProplessRenderer) SomeComponent {
	return &proplessWrapper{Renderer: renderer}
}

// Transformer for Propful -> SomeComponent

type propfulWrapper[TProps any] struct {
	BasicComponent

	Renderer Renderer[TProps]
	Props    TProps
}

func (c *propfulWrapper[TProps]) Render(width, height int) string {
	return c.Renderer(c.Props, width, height)
}

func PropfulToComponent[TProps any](renderer Renderer[TProps], props TProps) SomeComponent {
	return &propfulWrapper[TProps]{Renderer: renderer, Props: props}
}
