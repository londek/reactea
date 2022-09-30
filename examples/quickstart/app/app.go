package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"

	"github.com/londek/reactea/example/quickstart/pages/displayname"
	"github.com/londek/reactea/example/quickstart/pages/input"
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
	// Does it remind you of something? react-router!
	return c.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := input.New()

			return component, component.Init(input.Props{
				SetText: c.setText,
			})
		},
		"displayname": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := reactea.Componentify[string](displayname.Renderer)

			return component, component.Init(c.text)
		},
	})
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
