package render

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
		display:      Inline,
		width:        defaultLength,
		height:       defaultLength,
		maxWidth:     defaultLength,
		maxHeight:    defaultLength,
		minWidth:     defaultLength,
		minHeight:    defaultLength,
		margin:       zeroMargin,
		wrappingMode: ClipContent,
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

	wrappingMode WrappingMode

	renderer ResponsiveRenderer

	childrenContexts []*Context

	parent *Context
}

func (rc *Context) WrappingMode(wrappingMode WrappingMode) *Context {
	rc.wrappingMode = wrappingMode
	return rc
}

func (rc *Context) MinWidth(minWidth any) *Context {
	rc.minWidth = Length(minWidth)
	return rc
}

func (rc *Context) Width(width any) *Context {
	rc.width = Length(width)
	return rc
}

func (rc *Context) MaxWidth(maxWidth any) *Context {
	rc.maxWidth = Length(maxWidth)
	return rc
}

func (rc *Context) MinHeight(minHeight any) *Context {
	rc.minHeight = Length(minHeight)
	return rc
}

func (rc *Context) Height(height any) *Context {
	rc.height = Length(height)
	return rc
}

func (rc *Context) MaxHeight(maxHeight any) *Context {
	rc.maxHeight = Length(maxHeight)
	return rc
}

func (rc *Context) Margin(margin margin) *Context {
	rc.margin = margin
	return rc
}

// Shorthand for rc.MinWidth() and rc.MinHeight()
func (rc *Context) MinSize(minWidth, minHeight any) *Context {
	rc.minWidth = Length(minWidth)
	rc.minHeight = Length(minHeight)
	return rc
}

// Shorthand for rc.Width() and rc.Height()
func (rc *Context) Size(width, height any) *Context {
	rc.width = Length(width)
	rc.height = Length(height)
	return rc
}

// Shorthand for rc.MaxWidth() and rc.MaxHeight()
func (rc *Context) MaxSize(maxWidth, maxHeight any) *Context {
	rc.maxWidth = Length(maxWidth)
	rc.maxHeight = Length(maxHeight)
	return rc
}

func (rc *Context) Add(renderer Renderer) *Context {
	child := New()
	child.parent = rc

	renderer.Render(child)

	rc.childrenContexts = append(rc.childrenContexts, child)

	return child
}

// Creates container element and returns it
// Container behaves just like *render.Context
// but has Value() disabled
func (rc *Context) Container() container {
	child := New()
	child.parent = rc

	container := container{child}

	rc.childrenContexts = append(rc.childrenContexts, container.Context)

	return container
}

// Shorthand for rc.Add(render.Paragraph())
func (rc *Context) Paragraph(items ...any) *Context {
	return rc.Add(Paragraph(fmt.Sprint(items...)))
}

// Shorthand for rc.Add(render.Span())
func (rc *Context) Span(items ...any) *Context {
	return rc.Add(Span(fmt.Sprint(items...)))
}

// Shorthand for rc.Add(render.Breakline{})
func (rc *Context) Breakline() *Context {
	return rc.Add(Breakline{})
}

func (rc *Context) Render(availableWidthMaster, availableHeightMaster int) (result string) {
	if rc.contextType == RenderableKind {
		return rc.renderer(availableWidthMaster, availableHeightMaster)
	}

	var b strings.Builder
	b.Grow(availableWidthMaster * availableHeightMaster)

	availableWidthContext := availableWidthMaster
	availableHeightContext := availableHeightMaster

	breakLine := false
	for i, child := range rc.childrenContexts {
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

func (rc *Context) String() string {
	return fmt.Sprintf("display: %s | maxWidth: %s - width: %s - minWidth: %s | maxHeight: %s - height: %s - minWidth: %s | kind: %s | margin: %s", rc.display, rc.maxWidth, rc.width, rc.minWidth, rc.maxHeight, rc.height, rc.minWidth, rc.contextType, rc.margin)
}

func (rc *Context) Value(value string) *Context {
	return rc.Renderer(func(width, height int) string { return value })
}

func (rc *Context) Renderer(renderer ResponsiveRenderer) *Context {
	rc.contextType = RenderableKind
	rc.renderer = renderer
	return rc
}

func (rc *Context) TreeString(indent string) string {
	result := indent + "- " + rc.String() + "\n"

	for _, child := range rc.childrenContexts {
		result += child.TreeString("  " + indent)
	}

	return result
}

/*
Example:
func Render(rc *Context) {
	rc.Axis(Vertical)
	rc.Direction(MinToMax)

	rc.Add(c.footer)
	rc.Add(c.mainRouter)
	rc.Add(c.footer)
}

and next:
rc.String()
*/
