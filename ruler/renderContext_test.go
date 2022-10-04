package ruler

import (
	"fmt"
	"testing"
)

func TestRenderContext(t *testing.T) {
	Debug()

	rc := RenderContext{}

	rc.Add(Span("but nobody asks about how i feel"))
	rc.Add(Span("but nobody asks about how i feel"))
	rc.Add(Span("but nobody asks about how i feel"))

	rc.Add(BreakLine{})

	rc.Add(Span("but nobody asks about how i feel"))
	rc.Add(Span("but nobody asks about how i feel"))
	rc.Add(BreakLine{})

	rc.Add(Paragraph("but nobody asks about how i feel"))
	rc.Add(Paragraph("but nobody asks about how i feel"))

	fmt.Println(rc.String())
	fmt.Println(rc.TreeString(""))
}
