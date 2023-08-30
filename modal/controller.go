package modal

import (
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type Controller struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	initCmd tea.Cmd
	modal   reactea.SomeComponent
	mutex   sync.Mutex
}

func (c *Controller) Update(msg tea.Msg) tea.Cmd {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.initCmd != nil {
		initCmd := c.initCmd
		c.initCmd = nil
		return tea.Batch(initCmd, c.modal.Update(msg))
	}

	return c.modal.Update(msg)
}

func (c *Controller) Render(width, height int) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.modal.Render(width, height)
}

func (c *Controller) show(modal reactea.SomeComponent, initCmd tea.Cmd) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.modal = modal
	c.initCmd = initCmd
}

func (c *Controller) hide() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.modal.Destroy()
	c.modal = nil
}
