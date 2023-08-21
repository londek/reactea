package ruler

import (
	"fmt"
	"testing"
)

func TestContext(t *testing.T) {
	Debug()

	rc := New()

	rc.Span("abababab")
	rc.Add(Span("123365464574"))
	rc.Add(Span("foo bar foo baz"))

	rc.Add(Breakline{})

	rc.Add(Span("123365464574"))
	rc.Add(Span("foo bar foo baz"))
	rc.Breakline()

	container := rc.Container()
	container.Span("hello world")
	container.Width(5)
	container.Height(10)

	rc.Add(Span("foo bar foo baz"))
	rc.Add(Span("abababab"))

	fmt.Println(rc.Render(7, 50))
	fmt.Println(rc.TreeString(""))
}
