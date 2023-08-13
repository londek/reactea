package ruler

func Span(s string) span {
	return span(s)
}

type span string

func (s span) Render(rc *RenderContext) {
	rc.display = Inline

	rc.Value(string(s))
}
