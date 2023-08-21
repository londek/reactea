package render

type Breakline struct{}

func (br Breakline) Render(c *Context) {
	c.display = Inline

	c.Value("\n")
}
