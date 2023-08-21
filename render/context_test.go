package ruler

import (
	"fmt"
	"testing"
)

func TestContext(t *testing.T) {
	Debug()

	c := New()

	c.Span("abababab")
	c.Add(Span("123365464574"))
	c.Add(Span("foo bar foo baz"))

	c.Add(Breakline{})

	c.Add(Span("123365464574"))
	c.Add(Span("foo bar foo baz"))
	c.Breakline()

	container := c.Container()
	container.Span("hello world")
	container.Width(5)
	container.Height(10)

	c.Add(Span("foo bar foo baz"))
	c.Add(Span("abababab"))

	fmt.Println(c.Render(7, 50))
	fmt.Println(c.TreeString(""))
}
