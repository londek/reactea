package reactea

import tea "github.com/charmbracelet/bubbletea"

type RerenderMsg struct{}

// Utility tea.Cmd for requesting rerender
func Rerender() tea.Msg {
	return RerenderMsg{}
}

// Renders all AnyRenderers in one function
//
// Note: If you are using ProplessRenderer/DumbRenderer just pass
// reactea.NoProps{} or struct{}{}
//
// Note: Using named return type for 100% coverage
func RenderAny[TProps any, TRenderer AnyRenderer[TProps]](renderer TRenderer, props TProps, width, height int) (result string) {
	switch renderer := any(renderer).(type) {
	// TODO: Change to Renderer[TProps] along with
	// generics type-aliases feature (Planned Go 1.20)
	case func(TProps, int, int) string:
		result = renderer(props, width, height)
	case ProplessRenderer:
		result = renderer(width, height)
	case DumbRenderer:
		result = renderer()
	}

	return
}

// Handles rendering of all AnyProplessRenderers in one function
//
// Note: Using named return type for 100% coverage
func RenderDumb[TRenderer AnyProplessRenderer](renderer TRenderer, width, height int) (result string) {
	switch renderer := any(renderer).(type) {
	case ProplessRenderer:
		result = renderer(width, height)
	case DumbRenderer:
		result = renderer()
	}

	return
}

// Wraps propful into propless renderer
func PropfulToLess[TProps any](renderer Renderer[TProps], props TProps) ProplessRenderer {
	return func(width, height int) string {
		return renderer(props, width, height)
	}
}

// Static component for displaying static text

type staticComponent struct {
	BasicComponent

	content string
}

func (c *staticComponent) Render(int, int) string { return c.content }

func StaticComponent(content string) Component {
	return &staticComponent{content: content}
}

// Transformer for AnyRenderer -> Component

type componentTransformer[TProps any, TRenderer AnyRenderer[TProps]] struct {
	BasicComponent

	props    TProps
	renderer TRenderer
}

func (c *componentTransformer[TProps, TRenderer]) Render(width, height int) string {
	return RenderAny(c.renderer, c.props, width, height)
}

// Componentifies AnyRenderer
// Returns uninitialized component with renderer taking care of .Render()
func Componentify[TProps any, TRenderer AnyRenderer[TProps]](renderer TRenderer) Component {
	return &componentTransformer[TProps, TRenderer]{renderer: renderer}
}

// Transformer for AnyProplessRenderer -> Component

type dumbComponentTransformer[TRenderer AnyProplessRenderer] struct {
	BasicComponent

	renderer TRenderer
}

func (c *dumbComponentTransformer[T]) Render(width, height int) string {
	return RenderDumb(c.renderer, width, height)
}

// Componentifies AnyProplessRenderer
// Returns uninitialized component with renderer taking care of .Render()
func ComponentifyDumb[TRenderer AnyProplessRenderer](renderer TRenderer) Component {
	return &dumbComponentTransformer[TRenderer]{renderer: renderer}
}

// Used for tests
type mockComponent[TState any] struct {
	initFunc    func(Component, TState) tea.Cmd
	updateFunc  func(Component, TState) tea.Cmd
	renderFunc  func(Component, TState, int, int) string
	destroyFunc func(Component, TState)

	state TState
}

func (c *mockComponent[TState]) Init() tea.Cmd {
	return c.initFunc(c, c.state)
}

func (c *mockComponent[TState]) Update(msg tea.Msg) tea.Cmd {
	return c.updateFunc(c, c.state)
}

func (c *mockComponent[TState]) Render(width, height int) string {
	return c.renderFunc(c, c.state, width, height)
}

func (c *mockComponent[TState]) Destroy() {
	c.destroyFunc(c, c.state)
}
