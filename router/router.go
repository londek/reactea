package router

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type Params = map[string]string
type RouteInitializer func(Params) (reactea.SomeComponent, tea.Cmd)

type Component struct {
	reactea.BasicComponent[Props]

	lastComponent   reactea.SomeComponent
	lastPlaceholder string
}

type Props map[string]RouteInitializer

func New() *Component {
	return &Component{}
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)

	return c.initRoute()
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	if c.lastComponent == nil {
		return nil
	}

	return c.lastComponent.Update(msg)
}

func (c *Component) AfterUpdate() tea.Cmd {
	// If last route was changed we want to reuse the component
	if !reactea.WasRouteChanged() {
		return nil
	}

	// If last placeholder was wildcard and current route still matches
	// that wildcard, we want to reuse the component
	if _, ok := reactea.RouteMatchesPlaceholder(reactea.CurrentRoute(), c.lastPlaceholder); ok && c.lastPlaceholder != "" {
		return nil
	}

	if c.lastComponent != nil {
		c.lastComponent.Destroy()
	}

	c.lastComponent = nil

	return c.initRoute()
}

func (c *Component) Render(width, height int) string {
	if c.lastComponent != nil {
		return c.lastComponent.Render(width, height)
	}

	return fmt.Sprintf("Couldn't route for \"%s\"", reactea.CurrentRoute())
}

func (c *Component) initRoute() tea.Cmd {
	var cmd tea.Cmd

	if initializer, params, placeholder, ok := c.findMatchingRouteInitializer(); ok {
		c.lastComponent, cmd = initializer(params)
		c.lastPlaceholder = placeholder
	} else if initializer, ok := c.Props()["default"]; ok {
		c.lastComponent, cmd = initializer(nil)
	}

	return cmd
}

func (c *Component) findMatchingRouteInitializer() (RouteInitializer, Params, string, bool) {
	currentRoute := reactea.CurrentRoute()

	for placeholder, initializer := range c.Props() {
		if params, ok := reactea.RouteMatchesPlaceholder(currentRoute, placeholder); ok {
			return initializer, params, placeholder, true
		}
	}

	return nil, nil, "", false
}
