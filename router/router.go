package router

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type Params = map[string]string
type RouteInitializer func(Params) reactea.Component
type Routes = map[string]RouteInitializer

type Component struct {
	reactea.BasicComponent

	Routes Routes

	currentComponent reactea.Component
}

func New() *Component {
	return &Component{}
}

func NewWithRoutes(routes Routes) *Component {
	return &Component{Routes: routes}
}

func (c *Component) Init() tea.Cmd {
	return c.initRoute()
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	var initCmd, updateCmd tea.Cmd

	switch msg.(type) {
	case reactea.RouteUpdatedMsg:
		if c.currentComponent != nil {
			c.currentComponent.Destroy()
		}

		initCmd = c.initRoute()
	}

	if c.currentComponent != nil {
		updateCmd = c.currentComponent.Update(msg)
	}

	return tea.Batch(initCmd, updateCmd)
}

func (c *Component) Render(width, height int) string {
	if c.currentComponent != nil {
		return c.currentComponent.Render(width, height)
	}

	return fmt.Sprintf("Couldn't route for \"%s\"", reactea.CurrentRoute())
}

func (c *Component) initRoute() tea.Cmd {
	if initializer, params, ok := c.findMatchingRouteInitializer(); ok {
		c.currentComponent = initializer(params)
		return c.currentComponent.Init()
	}

	if initializer, ok := c.Routes["default"]; ok {
		c.currentComponent = initializer(nil)
		return c.currentComponent.Init()
	}

	return nil
}

func (c *Component) findMatchingRouteInitializer() (RouteInitializer, Params, bool) {
	currentRoute := reactea.CurrentRoute()

	for placeholder, initializer := range c.Routes {
		if params, ok := reactea.RouteMatchesPlaceholder(currentRoute, placeholder); ok {
			return initializer, params, true
		}
	}

	return nil, nil, false
}
