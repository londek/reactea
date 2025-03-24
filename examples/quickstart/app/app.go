package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"

	"github.com/londek/reactea/examples/quickstart/pages/displayname"
	"github.com/londek/reactea/examples/quickstart/pages/input"
	"github.com/londek/reactea/router"
)

type Component struct {
	reactea.BasicComponent

	mainRouter *router.Component
	text       string
}

func New() *Component {
	c := &Component{}

	c.mainRouter = router.NewWithRoutes(router.Routes{
		"default": func(router.Params) reactea.Component {
			component := input.New()

			component.SetText = c.setText

			return component
		},
		"/displayname": func(router.Params) reactea.Component {
			return reactea.Componentify(displayname.Render, c.text)
		},
	})

	return c
}

func (c *Component) Init() tea.Cmd {
	// Does it remind you of something? react-router!
	return c.mainRouter.Init()
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
	}

	return c.mainRouter.Update(msg)
}

func (c *Component) Render(width, height int) string {
	return c.mainRouter.Render(width, height)
}

func (c *Component) setText(text string) {
	c.text = text
}
