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
// making state mutable in any lifecycle method
//
// Note: Lifecycle is fully controlled by parent component
// making graph above fully theoretical and possibly
// invalid for third-party components

type Component interface {
	// Common lifecycle methods

	// Init() Is meant to both initialize subcomponents and run
	// long IO operations through tea.Cmd
	Init() tea.Cmd

	// It's called when component is about to be destroyed
	Destroy()

	// Typical tea.Model Update(), we handle all IO events here
	Update(tea.Msg) tea.Cmd

	// Render() is called when component should render itself
	// Provided width and height are target dimensions
	Render(int, int) string
}

// Why not Renderer[TProps]? It would have to be a type alias
// there are no type aliases yet for generics, but they are
// planned for some time soon. Something to keep in mind for future
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
// support is implemented. For now explicit
// type conversion is required
type Renderer[TProps any] func(TProps, int, int) string

// SUPEEEEEER shorthand for components
type ProplessRenderer = func(int, int) string

// Doesn't have state, props, even scalling for
// target dimensions = DumbRenderer, or Stringer
type DumbRenderer = func() string

// Alias for no props
type NoProps = struct{}

// Basic component that implements all methods
// required by reactea.Component
// except Render(int, int)
type BasicComponent struct{}

func (c *BasicComponent) Init() tea.Cmd              { return nil }
func (c *BasicComponent) Destroy()                   {}
func (c *BasicComponent) Update(msg tea.Msg) tea.Cmd { return nil }

// Utility component for displaying empty string on Render()
type InvisibleComponent struct{}

func (c *InvisibleComponent) Render(int, int) string { return "" }

// Destroys app before quiting
func Destroy() tea.Msg {
	return destroyAppMsg{}
}

type destroyAppMsg struct{}
