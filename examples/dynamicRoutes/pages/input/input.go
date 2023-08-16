package input

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/example/dynamicRoutes/data"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	textinput textinput.Model

	ids []int
}

func New() *Component {
	var ids []int
	for id := range data.Players {
		ids = append(ids, id)
	}

	return &Component{
		textinput: textinput.New(),
		ids:       ids,
	}
}

func (c *Component) Init(reactea.NoProps) tea.Cmd {
	return c.textinput.Focus()
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			// Validate input
			n, err := strconv.Atoi(c.textinput.Value())
			if err != nil {
				c.textinput.SetValue("Error")
				return nil
			}

			reactea.SetCurrentRoute(fmt.Sprintf("/players/%d", n))
			return nil
		}
	}

	var cmd tea.Cmd
	c.textinput, cmd = c.textinput.Update(msg)
	return cmd
}

func (c *Component) Render(int, int) string {
	return fmt.Sprintf("Found players with ids %v\nEnter player id: %s\nAnd press [ Enter ]", c.ids, c.textinput.View())
}
