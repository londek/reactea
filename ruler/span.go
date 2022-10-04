package ruler

type Span string

func (s Span) Render(rc *RenderContext) {
	rc.display = Inline

	rc.Value(string(s))
}
