package reactea

import tea "github.com/charmbracelet/bubbletea"

// It simplifies
// tea.NewProgram(reactea.New(root), opts...) to
// reactea.NewProgram(root, opts...)
//
// Note: Return type is *tea.Program, Reactea doesn't have
// it's own wrapper (reactea.Program) type, yet (?)
func NewProgram(root Component[NoProps], options ...tea.ProgramOption) *tea.Program {
	// Ensure globals are default, useful for tests and
	// running programs SEQUENTIALLY in same runtime
	isUpdate = false
	currentRoute = "/"
	lastRoute = "/"
	wasRouteChanged = false
	afterUpdaters = nil

	return tea.NewProgram(model{root, 0, 0}, options...)
}
