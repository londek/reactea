package ruler

type Paragraph string

func (p Paragraph) Render(rc *RenderContext) {
	rc.display = Block

	rc.Value(string(p))
}
