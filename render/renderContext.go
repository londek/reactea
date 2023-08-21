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

type Renderer interface{ Render(*renderContext) }

func New() *renderContext {
	return &renderContext{
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

// renderContext represents current level in Renderer context tree (reactea's Component tree)
// In TUI terms you could say renderContext is block
type renderContext struct {
	contextType Kind

	display display

	width, height       length
	maxWidth, maxHeight length
	minWidth, minHeight length

	margin margin

	renderer ResponsiveRenderer

	childrenContexts []*renderContext

	parent *renderContext
}

func (rc *renderContext) MinWidth(minWidth any) *renderContext {
	rc.minWidth = Length(minWidth)
	return rc
}

func (rc *renderContext) Width(width any) *renderContext {
	rc.width = Length(width)
	return rc
}

func (rc *renderContext) MaxWidth(maxWidth any) *renderContext {
	rc.maxWidth = Length(maxWidth)
	return rc
}

func (rc *renderContext) MinHeight(minHeight any) *renderContext {
	rc.minHeight = Length(minHeight)
	return rc
}

func (rc *renderContext) Height(height any) *renderContext {
	rc.height = Length(height)
	return rc
}

func (rc *renderContext) MaxHeight(maxHeight any) *renderContext {
	rc.maxHeight = Length(maxHeight)
	return rc
}

func (rc *renderContext) Margin(margin margin) *renderContext {
	rc.margin = margin
	return rc
}

// Shorthand for rc.MinWidth and rc.MinHeight
func (rc *renderContext) MinSize(minWidth, minHeight any) *renderContext {
	rc.minWidth = Length(minWidth)
	rc.minHeight = Length(minHeight)
	return rc
}

// Shorthand for rc.Width and rc.Height
func (rc *renderContext) Size(width, height any) *renderContext {
	rc.width = Length(width)
	rc.height = Length(height)
	return rc
}

// Shorthand for rc.MaxWidth and rc.MaxHeight
func (rc *renderContext) MaxSize(maxWidth, maxHeight any) *renderContext {
	rc.maxWidth = Length(maxWidth)
	rc.maxHeight = Length(maxHeight)
	return rc
}

func (rc *renderContext) Add(renderer Renderer) *renderContext {
	child := New()
	child.parent = rc

	renderer.Render(child)

	rc.childrenContexts = append(rc.childrenContexts, child)

	return child
}

func (rc *renderContext) AddContainer() container {
	child := New()
	child.parent = rc

	container := container{child}

	rc.childrenContexts = append(rc.childrenContexts, container.renderContext)

	return container
}

func (rc *renderContext) Render(availableWidthMaster, availableHeightMaster int) (result string) {
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

func (rc *renderContext) String() string {
	return fmt.Sprintf("display: %s | maxWidth: %s - width: %s - minWidth: %s | maxHeight: %s - height: %s - minWidth: %s | kind: %s | margin: %s", rc.display, rc.maxWidth, rc.width, rc.minWidth, rc.maxHeight, rc.height, rc.minWidth, rc.contextType, rc.margin)
}

func (rc *renderContext) Value(value string) *renderContext {
	return rc.Renderer(func(width, height int) string { return value })
}

func (rc *renderContext) Renderer(renderer ResponsiveRenderer) *renderContext {
	rc.contextType = RenderableKind
	rc.renderer = renderer
	return rc
}

func (rc *renderContext) TreeString(indent string) string {
	result := indent + "- " + rc.String() + "\n"

	for _, child := range rc.childrenContexts {
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
