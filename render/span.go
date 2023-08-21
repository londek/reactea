package render

type Span string

func (s Span) Render(rc *Context) {
	rc.display = Inline

	rc.Value(string(s))
}
