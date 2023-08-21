package ruler

type Span string

func (s Span) Render(rc *renderContext) {
	rc.display = Inline

	rc.Value(string(s))
}
