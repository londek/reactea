package reactea

import tea "github.com/charmbracelet/bubbletea"

// The lifecycle is
//
//           \/ Usually won't be called on first render
// Init ---> Update -> Render ---> Destroy?
//       |                     |   /\ implementation detail and
//       |---------------------|   therefore doesn't return tea.Cmd
//
// Reactea takes pointer approach for components
// making state modifiable in any lifecycle method
//
// Note: Remember that Update IS NOT guaranteed to be called on
// first-run, Init() is, and critical logic should be there
//
// Note: Lifecycle is fully controlled by parent component
// making graph above fully theoretical and possibly
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
	// Note: It's parent component's job to call it so
	// relying on it outside of Reactea builtins is
	// not reliable
	Destroy()

	// Typical tea.Model Update(), we handle all IO events here
	Update(tea.Msg) tea.Cmd

	// Callee already knows at which size should it render at
	Render(int, int) string

	// Reactea implementation methods, just use BasicComponent
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
// without the props part, use it if you don't care about generics
type SomeComponent interface {
	Destroy()
	Update(tea.Msg) tea.Cmd
	Render(int, int) string
	AfterUpdate() tea.Cmd
}

// Why not Renderer[TProps]? It would have to be type alias
// there are no type aliases yet for generics, but they are
// planned for Go 1.20. Something to keep in mind for future
type AnyRenderer[TProps any] interface {
	func(TProps, int, int) string | AnyProplessRenderer
}

type AnyProplessRenderer interface {
	ProplessRenderer | DumbRenderer
}

// Ultra shorthand for components = just renderer
// One could say it's a stateless component
// Also note that it doesn't handle any IO by itself
//
// TODO: Change to type alias after type aliases for generics
// support is released (planned Go 1.20). For now explicit
// type conversion is required
type Renderer[TProps any] func(TProps, int, int) string

// SUPEEEEEER shorthand for components
type ProplessRenderer = func(int, int) string

// Doesn't have state, props, even scalling for
// target dimensions = DumbRenderer, or Stringer
type DumbRenderer = func() string

// Basic component that implements all methods
// required by reactea.Component[TProps]
// except Render(int, int)
type BasicComponent[TProps any] struct {
	props TProps
}

func (c *BasicComponent[TProps]) Init(props TProps) tea.Cmd  { c.UpdateProps(props); return nil }
func (c *BasicComponent[TProps]) Destroy()                   {}
func (c *BasicComponent[TProps]) Update(msg tea.Msg) tea.Cmd { return nil }
func (c *BasicComponent[TProps]) AfterUpdate() tea.Cmd       { return nil }
func (c *BasicComponent[TProps]) UpdateProps(props TProps)   { c.props = props }
func (c *BasicComponent[TProps]) Props() TProps              { return c.props }

// The most basic form of Reactea Component
// It implements all not required methods
// so you don't have to
type BasicProplessComponent struct{}

func (c *BasicProplessComponent) Destroy()                   {}
func (c *BasicProplessComponent) Update(msg tea.Msg) tea.Cmd { return nil }
func (c *BasicProplessComponent) AfterUpdate() tea.Cmd       { return nil }

// Utility component for displaying empty string on Render()
type InvisibleComponent struct{}

func (c *InvisibleComponent) Render(int, int) string { return "" }

type NoProps = struct{}

// Destroys app before quiting
func Destroy() tea.Msg {
	return destroyAppMsg{}
}

type destroyAppMsg struct{}
