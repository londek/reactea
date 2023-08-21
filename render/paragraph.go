package render

type Paragraph string

func (p Paragraph) Render(rc *Context) {
	rc.display = Block

	rc.Value(string(p))
}
