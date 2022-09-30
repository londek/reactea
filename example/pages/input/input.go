package input

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	textinput textinput.Model
}

type Props struct {
	SetText func(string)
}

func New() *Component {
	return &Component{
		textinput: textinput.New(),
	}
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)

	return c.textinput.Focus()
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			// Lifted state power! Woohooo
			c.Props().SetText(c.textinput.Value())

			reactea.SetCurrentRoute("displayname")

			return nil
		}
	}

	var cmd tea.Cmd
	c.textinput, cmd = c.textinput.Update(msg)
	return cmd
}

// Here we are not using width and height, but you can!
// Using lipgloss styles for example
func (c *Component) Render(int, int) string {
	return fmt.Sprintf("Enter your name: %s\nAnd press [ Enter ]", c.textinput.View())
}
