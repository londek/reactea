package ruler

import (
	"fmt"
	"strings"
)

var debug bool

func Debug() {
	debug = true
}

type RenderedElement struct {
	inline  bool
	element string
}

type RenderedElements []RenderedElement

func (re RenderedElements) Join() (result string) {
	result += re[0].element

	for _, element := range re[1:] {
		if !element.inline {
			result += "\n"
		}

		result += element.element
	}

	return
}

type ResponsiveRenderer = func(width, height int) string

type Renderer interface{ Render(*Context) }

func New() *Context {
	return &Context{
		display:   Inline,
		width:     defaultLength,
		height:    defaultLength,
		maxWidth:  defaultLength,
		maxHeight: defaultLength,
		minWidth:  defaultLength,
		minHeight: defaultLength,
		margin:    zeroMargin,
	}
}

// So that we dont have to recalculate everything
type cachedElement struct {
	element       string
	width, height int
}

// Context represents current level in Renderer context tree (reactea's Component tree)
// In TUI terms you could say Context is block
type Context struct {
	contextType Kind

	display display

	width, height       length
	maxWidth, maxHeight length
	minWidth, minHeight length

	margin margin

	renderer ResponsiveRenderer

	childrenContexts []*Context

	parent *Context
}

func (c *Context) MinWidth(minWidth any) *Context {
	c.minWidth = Length(minWidth)
	return c
}

func (c *Context) Width(width any) *Context {
	c.width = Length(width)
	return c
}

func (c *Context) MaxWidth(maxWidth any) *Context {
	c.maxWidth = Length(maxWidth)
	return c
}

func (c *Context) MinHeight(minHeight any) *Context {
	c.minHeight = Length(minHeight)
	return c
}

func (c *Context) Height(height any) *Context {
	c.height = Length(height)
	return c
}

func (c *Context) MaxHeight(maxHeight any) *Context {
	c.maxHeight = Length(maxHeight)
	return c
}

func (c *Context) Margin(margin margin) *Context {
	c.margin = margin
	return c
}

// Shorthand for c.MinWidth and c.MinHeight
func (c *Context) MinSize(minWidth, minHeight any) *Context {
	c.minWidth = Length(minWidth)
	c.minHeight = Length(minHeight)
	return c
}

// Shorthand for c.Width and c.Height
func (c *Context) Size(width, height any) *Context {
	c.width = Length(width)
	c.height = Length(height)
	return c
}

// Shorthand for c.MaxWidth and c.MaxHeight
func (c *Context) MaxSize(maxWidth, maxHeight any) *Context {
	c.maxWidth = Length(maxWidth)
	c.maxHeight = Length(maxHeight)
	return c
}

func (c *Context) Add(renderer Renderer) *Context {
	child := New()
	child.parent = c

	renderer.Render(child)

	c.childrenContexts = append(c.childrenContexts, child)

	return child
}

// Creates container element and returns it
// Container behaves just like *render.Context
// but has Value() disabled
func (c *Context) Container() container {
	child := New()
	child.parent = c

	container := container{child}

	c.childrenContexts = append(c.childrenContexts, container.Context)

	return container
}

// Shorthand for c.Add(render.Paragraph())
func (c *Context) Paragraph(items ...any) *Context {
	return c.Add(Paragraph(fmt.Sprint(items...)))
}

// Shorthand for c.Add(render.Span())
func (c *Context) Span(items ...any) *Context {
	return c.Add(Span(fmt.Sprint(items...)))
}

// Shorthand for c.Add(render.Breakline{})
func (c *Context) Breakline() *Context {
	return c.Add(Breakline{})
}

func (c *Context) Render(availableWidthMaster, availableHeightMaster int) (result string) {
	if c.contextType == RenderableKind {
		return c.renderer(availableWidthMaster, availableHeightMaster)
	}

	var b strings.Builder
	b.Grow(availableWidthMaster * availableHeightMaster)

	availableWidthContext := availableWidthMaster
	availableHeightContext := availableHeightMaster

	breakLine := false
	for i, child := range c.childrenContexts {
		var rendered string

		switch child.display {
		// dimensions, margins are fully ignored by inline elements
		case Inline:
			if breakLine {
				b.WriteRune('\n')
				breakLine = false
			}

			rendered = child.Render(200, 50)
			width, height := SizeOf(rendered)
			availableWidthContext -= width
			availableHeightContext -= height

			if availableWidthContext < 0 {
				breakLine = true
				availableWidthContext = availableWidthMaster
			}

			if width > availableWidthMaster {
				rendered = ClipWidth(rendered, availableWidthMaster)
			}

		// blocks are breakline elements, they reset availableWidth
		case Block:
			if i > 0 {
				b.WriteRune('\n')
			}

			rendered = child.Render(availableWidthMaster, availableHeightContext)
			width, height := SizeOf(rendered)
			availableHeightContext -= height

			if width > availableWidthMaster {
				rendered = ClipWidth(rendered, availableWidthMaster)
			}

			availableWidthContext = availableWidthMaster
			breakLine = true
		case InlineBlock:
			if breakLine {
				b.WriteRune('\n')
				breakLine = false
			}

		}

		b.WriteString(rendered)

		if availableHeightContext < 0 {
			return ClipHeight(b.String(), availableHeightContext)
		} else if availableHeightContext == 0 {
			return b.String()
		}
	}

	return b.String()
}

func (c *Context) String() string {
	return fmt.Sprintf("display: %s | maxWidth: %s - width: %s - minWidth: %s | maxHeight: %s - height: %s - minWidth: %s | kind: %s | margin: %s", c.display, c.maxWidth, c.width, c.minWidth, c.maxHeight, c.height, c.minWidth, c.contextType, c.margin)
}

func (c *Context) Value(value string) *Context {
	return c.Renderer(func(width, height int) string { return value })
}

func (c *Context) Renderer(renderer ResponsiveRenderer) *Context {
	c.contextType = RenderableKind
	c.renderer = renderer
	return c
}

func (c *Context) TreeString(indent string) string {
	result := indent + "- " + c.String() + "\n"

	for _, child := range c.childrenContexts {
		result += child.TreeString("  " + indent)
	}

	return result
}

/*
Example:
func Render(rc RenderContext) {
	rc.Axis(Vertical)
	rc.Direction(MinToMax)

	rc.Add(c.footer)
	rc.Add(c.mainRouter)
	rc.Add(c.footer)
}

and next:
rc.String()
*/
