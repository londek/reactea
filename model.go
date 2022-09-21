package reactea

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Useful for constraining some actions to update-stage only
var isUpdate bool

type model struct {
	root Component[NoProps]

	width, height int
}

func New(root Component[NoProps]) model {
	return model{
		root: root,
	}
}

func (m model) Init() tea.Cmd {
	return m.root.Init(NoProps{})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	isUpdate = true

	switch msg := msg.(type) {
	case destroyAppMsg:
		m.root.Destroy()
		return m, tea.Quit
	// We want component to know at what size should it render
	// and unify size handling across all reactea components
	// That's why we don't forward WindowSizeMsg, but tell
	// target width and height with Render()
	//
	// Note: It also means no Update() before first render
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	}

	var cmd tea.Cmd
	cmd = m.root.Update(msg)

	isUpdate = false

	return m, tea.Sequence(cmd, handleAfterUpdates())
}

func (m model) View() string {
	return m.root.Render(m.width, m.height)
}
