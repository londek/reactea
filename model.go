package reactea

import tea "github.com/charmbracelet/bubbletea"

// Useful for constraining some actions to update-stage only
var isUpdate bool

type model struct {
	root Component

	width, height int
}

func New(root Component) model {
	return model{
		root: root,
	}
}

func (m model) Init() tea.Cmd {
	return m.root.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	handleBeforeUpdates()

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

	rootCmd := m.root.Update(msg)

	isUpdate = false

	afterUpdatesCmd := handleAfterUpdates()

	// Guarantee rerender if route was changed
	if wasRouteChanged {
		return m, tea.Batch(updatedRoute(lastRoute), rootCmd, afterUpdatesCmd)
	}

	return m, tea.Batch(rootCmd, afterUpdatesCmd)
}

func (m model) View() string {
	return m.root.Render(m.width, m.height)
}
