package ruler

import (
	"fmt"
	"testing"
)

func TestRenderContext(t *testing.T) {
	Debug()

	rc := &RenderContext{}

	rc.axis = Vertical
	rc.direction = MaxToMin

	rc.width = 30
	rc.height = 15

	rc.AddParagraph("hej")
	rc.AddParagraph("co")
	rc.AddParagraph("tam")
	rc.AddParagraph("u was")
	rc.AddParagraph("ta?")
	rc.Add(func(rc *RenderContext) {
		rc.width = 20
		rc.height = 2

		rc.direction = MaxToMin

		rc.AddParagraph("hej")
		rc.AddParagraph("roman!")
	})

	rc.AddParagraph("but nobody asks about how i feel")

	fmt.Println(rc.String())
}
