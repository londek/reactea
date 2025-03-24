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

	mainRouter reactea.Component

	text string
}

func New() *Component {
	return &Component{
		mainRouter: router.NewWithRoutes(map[string]router.RouteInitializer{
			"default": func(router.Params) reactea.Component {
				component := input.New()

				return component
			},
			// We are using dynamic routes (route params) in this example
			"/players/:playerId": func(params router.Params) reactea.Component {
				playerId, _ := strconv.Atoi(params["playerId"])

				component := reactea.Componentify(displayplayer.Render, playerId)

				return component
			},
		}),
	}
}

func (c *Component) Init(reactea.NoProps) tea.Cmd {
	return c.mainRouter.Init()
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return reactea.Destroy
		case "u":
			reactea.SetRoute("/")
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
