package modal

import (
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type Controller struct {
	reactea.BasicComponent

	initCmd tea.Cmd
	modal   reactea.Component
	mutex   sync.Mutex
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Update(msg tea.Msg) tea.Cmd {
	c.mutex.Lock()

	if c.initCmd != nil {
		initCmd := c.initCmd
		c.initCmd = nil
		return tea.Sequence(initCmd, c.modal.Update(msg))
	}

	return c.modal.Update(msg)
}

func (c *Controller) Render(width, height int) string {
	defer c.mutex.Unlock()

	return c.modal.Render(width, height)
}

func (c *Controller) Run(f func(*Controller) tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		return f(c)
	}
}

func (c *Controller) show(modal reactea.Component, initCmd tea.Cmd) {
	c.modal = modal
	c.initCmd = initCmd
}

func (c *Controller) hide() {
	c.modal.Destroy()
	c.modal = nil
}
