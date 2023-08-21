package ruler

type Paragraph string

func (p Paragraph) Render(rc *renderContext) {
	rc.display = Block

	rc.Value(string(p))
}
