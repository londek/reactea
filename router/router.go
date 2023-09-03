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
	initCmd         tea.Cmd
	lastPlaceholder string
}

func New() *Component {
	return &Component{}
}

func NewWithRoutes(routes map[string]RouteInitializer) *Component {
	return &Component{Routes: routes}
}

func (c *Component) Init() tea.Cmd {
	return c.initRoute()
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.BeforeUpdate(c)

	if c.lastComponent == nil {
		return nil
	}

	if c.initCmd != nil {
		initCmd := c.initCmd
		c.initCmd = nil
		return tea.Sequence(initCmd, c.lastComponent.Update(msg))
	}

	return c.lastComponent.Update(msg)
}

func (c *Component) BeforeUpdate() {
	// If last route was changed we want to reuse the component
	if !reactea.WasRouteChanged() {
		return
	}

	// If last placeholder was wildcard and current route still matches
	// that wildcard, we want to reuse the component
	if _, ok := reactea.RouteMatchesPlaceholder(reactea.CurrentRoute(), c.lastPlaceholder); ok && c.lastPlaceholder != "" {
		return
	}

	if c.lastComponent != nil {
		c.lastComponent.Destroy()
	}

	c.lastComponent = nil

	c.initRoute()
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
		c.initCmd = c.lastComponent.Init()
	} else if initializer, ok := c.Routes["default"]; ok {
		c.lastPlaceholder = placeholder
		c.lastComponent = initializer(nil)
		c.initCmd = c.lastComponent.Init()
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
