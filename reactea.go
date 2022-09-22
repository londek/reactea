package reactea

import (
	tea "github.com/charmbracelet/bubbletea"
)

// The lifecycle is
//
//           \/ Usually won't be called on first render
// Init ---> Update -> Render ---> Destroy?
//       |                     |   /\ implementation detail and
//       |---------------------|   therefore doesn't return tea.Cmd
//
// reactea takes pointer approach for components
// making state modifiable in any lifecycle method
//
// Note: Remember that Update IS NOT guaranteed to be called on
// first-run, Init() is, and critical logic should be there
//
// Note: Lifecycle is fully controlled by parent component
// making graph above fully theoritical and possibly
// invalid for third-party components

type Component[TProps any] interface {
	// Common lifecycle methods

	// You always have to initialize component with some kind of
	// props - it can even be zero value
	// Init() Is meant to both initialize subcomponents and run
	// long IO operations through tea.Cmd
	Init(TProps) tea.Cmd

	// It's called when component is about to be destroyed
	//
	// Note: It's parent component job to call it and
	// relying on it outside of reactea builtins is
	// not reliable
	Destroy()

	// Typical tea.Model Update(), we handle all IO events here
	Update(tea.Msg) tea.Cmd

	// Callee already knows at which size should it render at
	Render(int, int) string

	// reactea implementation methods, just use BasicComponent
	// if you don't know what are you doing

	// It's an Update() but for props, the state derive stage
	// happens here
	UpdateProps(TProps)

	// AfterUpdate is stage useful for components like routers
	// to prepare content. Saying that you will probably never
	// need to use it
	AfterUpdate() tea.Cmd
}

// Interface which is basically reactea.Component just
// without the props part
type SomeComponent interface {
	Destroy()
	Update(tea.Msg) tea.Cmd
	Render(int, int) string
	AfterUpdate() tea.Cmd
}

// I decided to give it a name "Component" and not "Renderer"
// Because Component.Render is itself ProplessRenderer
// (or the other way around ProplessRenderer = Component.Render)
// So the naming is infact valid for any type of
// components in reactea
type AnyComponent[TProps any] interface {
	Renderer[TProps] | AnyProplessComponent
}

// I decided to give it a name "Component" and not "Renderer"
// Because Component.Render is itself ProplessRenderer
// (or the other way around ProplessRenderer = Component.Render)
// So the naming is infact valid for any type of
// components in reactea
type AnyProplessComponent interface {
	ProplessRenderer | DumbRenderer
}

// Ultra shorthand for components = just renderer
// One could say it's a stateless component
// Also note that it doesn't handle any IO by itself
type Renderer[TProps any] func(TProps, int, int) string

// SUPEEEEEER shorthand for components
type ProplessRenderer func(int, int) string

// Doesn't have state, props, even scalling for
// target dimensions = DumbRenderer, or Stringer
type DumbRenderer func() string

// The most basic form of reactea component
// It implements all not required methods
// so you don't have to
type BasicComponent struct{}

func (c *BasicComponent) Destroy()                   {}
func (c *BasicComponent) Update(msg tea.Msg) tea.Cmd { return nil }
func (c *BasicComponent) AfterUpdate() tea.Cmd       { return nil }

// Stores props in struct.
// If you want to derive state from UpdateProps()
// you're probably looking at wrong thing
type BasicPropfulComponent[TProps any] struct {
	props TProps
}

func (c *BasicPropfulComponent[TProps]) Init(props TProps) tea.Cmd { c.UpdateProps(props); return nil }
func (c *BasicPropfulComponent[TProps]) UpdateProps(props TProps)  { c.props = props }
func (c *BasicPropfulComponent[TProps]) Props() TProps             { return c.props }

// Utility component for displaying empty string on Render()
type InvisibleComponent struct{}

func (c *InvisibleComponent) Render(int, int) string { return "" }

type NoProps struct{}

// Destroys app before quiting
func Destroy() tea.Msg {
	return destroyAppMsg{}
}

type destroyAppMsg struct{}
