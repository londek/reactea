package reactea

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Useful for constraining some actions to update-stage only
var isUpdate bool

type model struct {
	program *tea.Program
	root    Component

	width, height int
}

func (m model) Init() tea.Cmd {
	return m.root.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.execute(handleBeforeUpdates())

	isUpdate = true
	wasRouteChanged = false

	switch msg := msg.(type) {
	case destroyAppMsg:
		m.root.Destroy()
		return m, tea.Quit
	// We want component to know at what size should it render
	// and unify size handling across all Reactea components
	// We pass WindowSizeMsg to root component just for
	// sake of utility.
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	}

	m.execute(m.root.Update(msg))

	isUpdate = false

	m.execute(handleAfterUpdates())

	// Guarantee rerender if route was changed
	if wasRouteChanged {
		return m, updatedRoute(lastRoute)
	}

	return m, nil
}

func (m model) View() string {
	return m.root.Render(m.width, m.height)
}

func (m model) execute(cmd tea.Cmd) {
	if cmd == nil {
		return
	}

	go func() {
		msg := cmd()
		switch msg := msg.(type) {
		case tea.BatchMsg:
			for _, cmd := range msg {
				go func(cmd tea.Cmd) {
					m.program.Send(cmd())
				}(cmd)
			}
		default:
			m.program.Send(msg)
		}
	}()
}
