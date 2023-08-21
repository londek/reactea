package ruler

type BreakLine struct{}

func (br BreakLine) Render(rc *renderContext) {
	rc.display = Inline

	rc.Value("\n")
}
