package app

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"

	"github.com/londek/reactea/example/dynamicRoutes/pages/displayplayer"
	"github.com/londek/reactea/example/dynamicRoutes/pages/input"
	"github.com/londek/reactea/router"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter reactea.Component[router.Props]

	text string
}

func New() *Component {
	return &Component{
		mainRouter: router.New(),
	}
}

func (c *Component) Init(reactea.NoProps) tea.Cmd {
	return c.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := input.New()

			return component, component.Init(reactea.NoProps{})
		},
		// We are using dynamic routes (route params) in this example
		"players/:playerId": func(params router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := reactea.Componentify[int](displayplayer.Render)

			playerId, _ := strconv.Atoi(params["playerId"])

			return component, component.Init(playerId)
		},
	})
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return reactea.Destroy
		case "u":
			reactea.SetCurrentRoute("")
		}
	}

	return c.mainRouter.Update(msg)
}

func (c *Component) Render(width, height int) string {
	return fmt.Sprintf("Current route: \"%s\"\n\n%s", reactea.CurrentRoute(), c.mainRouter.Render(width, height))
}

func (c *Component) setText(text string) {
	c.text = text
}
