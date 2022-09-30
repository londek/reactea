package reactea

// Renders all AnyRenderers in one function
//
// Note: If you are using ProplessRenderer/DumbRenderer just pass
// reactea.NoProps{} or struct{}{}
func RenderAny[TProps any, TRenderer AnyRenderer[TProps]](renderer TRenderer, props TProps, width, height int) string {
	switch renderer := any(renderer).(type) {
	// TODO: Change to Renderer[TProps] along with
	// generics type-aliases feature (Planned Go 1.20)
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
