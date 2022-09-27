package reactea

// Renders all AnyRenderers in one function
func RenderAny[TProps any, TRenderer AnyRenderer[TProps]](renderer TRenderer, props TProps, width, height int) string {
	switch renderer := any(renderer).(type) {
	// TODO: Change to Renderer[TProps] along with type-aliases
	// for generics feature (Planned Go 1.20)
	case func(TProps, int, int) string:
		return renderer(props, width, height)
	case ProplessRenderer:
		return renderer(width, height)
	case DumbRenderer:
		return renderer()
	}

	return ""
}

// Wraps propful into propless renderer
func PropfulToLess[TProps any](renderer Renderer[TProps], props TProps) ProplessRenderer {
	return func(width, height int) string {
		return renderer(props, width, height)
	}
}

// Transformer for AnyRenderer -> Component

type componentTransformer[TProps any, TRenderer AnyRenderer[TProps]] struct {
	BasicComponent
	BasicPropfulComponent[TProps]

	renderer TRenderer
}

func (c *componentTransformer[TProps, TRenderer]) Render(width, height int) string {
	return RenderAny(c.renderer, c.props, width, height)
}

// Componentifies AnyRenderer
// Returns uninitialized component with renderer taking care of .Render()
func Componentify[TProps any, TRenderer AnyRenderer[TProps]](renderer TRenderer) Component[TProps] {
	return &componentTransformer[TProps, TRenderer]{renderer: renderer}
}

// Transformer for AnyRenderer -> Component

type someComponentTransformer[TProps any, TRenderer AnyRenderer[TProps]] struct {
	BasicComponent

	renderer TRenderer
	props    TProps
}

func (c *someComponentTransformer[TProps, TRenderer]) Render(width, height int) string {
	return RenderAny(c.renderer, c.props, width, height)
}

// I don't care about props, just give me SomeComponent
func SomeComponentify[TProps any, TRenderer AnyRenderer[TProps]](renderer TRenderer, initialProps TProps) SomeComponent {
	return &someComponentTransformer[TProps, TRenderer]{
		renderer: renderer,
		props:    initialProps,
	}
}
