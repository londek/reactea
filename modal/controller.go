package modal

import (
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type Controller struct {
	reactea.BasicComponent

	initFunc   func(*Controller) func() tea.Cmd
	escapeFunc func() tea.Cmd // Called when modal flow is ended

	rendered       string
	initCmd        tea.Cmd
	shouldDestruct bool

	modal    reactea.Component
	locked   bool
	cond     *sync.Cond
	w        waiter
	runMutex sync.Mutex
}

func NewController(initFunc func(*Controller) func() tea.Cmd) *Controller {
	return &Controller{
		initFunc: initFunc,
		cond:     sync.NewCond(&sync.Mutex{}),
		w:        make(waiter),
	}
}

func (c *Controller) Init() tea.Cmd {
	return c.Run(c.initFunc)
}

func (c *Controller) Update(msg tea.Msg) tea.Cmd {
	var initCmd, updateCmd, escapeCmd tea.Cmd

	c.cond.L.Lock()
	c.locked = true
	for c.modal == nil {
		c.cond.Wait()
	}

	updateCmd = c.modal.Update(msg)

	if c.escapeFunc != nil {
		escapeFunc := c.escapeFunc
		c.escapeFunc = nil
		escapeCmd = escapeFunc()
	}

	if c.initCmd != nil {
		initCmd = c.initCmd
		c.initCmd = nil
	}

	return tea.Batch(initCmd, updateCmd, escapeCmd)
}

func (c *Controller) Render(width, height int) string {
	if c.locked {
		defer c.cond.L.Unlock()
		c.locked = false
	}

	c.rendered = c.modal.Render(width, height)

	if c.shouldDestruct {
		c.modal.Destroy()
		c.modal = nil
		c.shouldDestruct = false
	}

	return c.rendered
}

func (c *Controller) Run(f func(*Controller) func() tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		c.runMutex.Lock()
		defer c.runMutex.Unlock()

		c.escapeFunc = f(c)
		c.w.Signal()

		return nil
	}
}
