package router

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type RouteInitializer func() (reactea.SomeComponent, tea.Cmd)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	// If last route != currentRoute we want to reinitialize the component
	lastRoute     reactea.Route
	lastComponent reactea.SomeComponent
}

type Props map[string]RouteInitializer

func New() *Component {
	return &Component{}
}

func (c *Component) Init(props Props) (cmd tea.Cmd) {
	c.UpdateProps(props)

	if initializer, ok := c.Props()["default"]; ok {
		c.lastComponent, cmd = initializer()
	}

	return
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	if c.lastComponent == nil {
		return nil
	}

	return c.lastComponent.Update(msg)
}

func (c *Component) AfterUpdate() tea.Cmd {
	defer func() { c.lastRoute = reactea.CurrentRoute() }()

	if reactea.CurrentRoute().Equal(c.lastRoute) {
		return nil
	}

	if c.lastComponent != nil {
		c.lastComponent.Destroy()
		c.lastComponent = nil
	}

	var cmd tea.Cmd

	if initializer, ok := c.Props()[reactea.CurrentRoute().String()]; ok {
		c.lastComponent, cmd = initializer()
	} else if initializer, ok := c.Props()["default"]; ok {
		c.lastComponent, cmd = initializer()
	}

	return cmd
}

func (c *Component) Render(width, height int) string {
	if c.lastComponent != nil {
		return c.lastComponent.Render(width, height)
	}

	return fmt.Sprintf("Couldn't route for '%s'", reactea.CurrentRoute())
}
