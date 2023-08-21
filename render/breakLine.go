package ruler

type Breakline struct{}

func (br Breakline) Render(rc *Context) {
	rc.display = Inline

	rc.Value("\n")
}
