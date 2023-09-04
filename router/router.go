package router

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type Params = map[string]string
type RouteInitializer func(Params) reactea.Component

type Component struct {
	reactea.BasicComponent

	Routes map[string]RouteInitializer

	lastComponent   reactea.Component
	lastPlaceholder string
}

func New() *Component {
	return &Component{}
}

func NewWithRoutes(routes map[string]RouteInitializer) *Component {
	return &Component{Routes: routes}
}

func (c *Component) Init() tea.Cmd {
	reactea.BeforeUpdate(c)
	return c.initRoute()
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.BeforeUpdate(c)

	if c.lastComponent == nil {
		return nil
	}

	return c.lastComponent.Update(msg)
}

func (c *Component) BeforeUpdate() tea.Cmd {
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
	if initializer, params, placeholder, ok := c.findMatchingRouteInitializer(); ok {
		c.lastPlaceholder = placeholder
		c.lastComponent = initializer(params)
		return c.lastComponent.Init()
	} else if initializer, ok := c.Routes["default"]; ok {
		c.lastPlaceholder = placeholder
		c.lastComponent = initializer(nil)
		return c.lastComponent.Init()
	}

	return nil
}

func (c *Component) findMatchingRouteInitializer() (RouteInitializer, Params, string, bool) {
	currentRoute := reactea.CurrentRoute()

	for placeholder, initializer := range c.Routes {
		if params, ok := reactea.RouteMatchesPlaceholder(currentRoute, placeholder); ok {
			return initializer, params, placeholder, true
		}
	}

	return nil, nil, "", false
}
