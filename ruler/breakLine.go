package ruler

type BreakLine struct{}

func (br BreakLine) Render(rc *RenderContext) {
	rc.display = Inline

	rc.Value("\n")
}
