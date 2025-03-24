package app

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"

	"github.com/londek/reactea/examples/dynamicRoutes/pages/displayplayer"
	"github.com/londek/reactea/examples/dynamicRoutes/pages/input"
	"github.com/londek/reactea/router"
)

type Component struct {
	reactea.BasicComponent

	mainRouter reactea.Component
}

func New() *Component {
	return &Component{
		mainRouter: router.NewWithRoutes(map[string]router.RouteInitializer{
			"default": func(router.Params) reactea.Component {
				return input.New()
			},
			// We are using dynamic routes (route params) in this example
			"/players/:playerId": func(params router.Params) reactea.Component {
				playerId, _ := strconv.Atoi(params["playerId"])

				return reactea.Componentify(displayplayer.Render, playerId)
			},
		}),
	}
}

func (c *Component) Init() tea.Cmd {
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
