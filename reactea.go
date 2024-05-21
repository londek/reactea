package reactea

import tea "github.com/charmbracelet/bubbletea"

// Note: Return type is *tea.Program, Reactea doesn't have
// it's own wrapper (reactea.Program) type, yet (?)
func NewProgram(root Component, options ...tea.ProgramOption) *tea.Program {
	// Ensure globals are default, useful for tests and
	// running programs SEQUENTIALLY during runtime
	isUpdate = false
	currentRoute = "/"
	lastRoute = "/"
	wasRouteChanged = false

	m := &model{
		program: nil,
		root:    root,
		width:   0,
		height:  0,
	}

	program := tea.NewProgram(m, options...)
	m.program = program

	return program
}
