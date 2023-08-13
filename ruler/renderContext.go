package ruler

import (
	"github.com/AvraamMavridis/randomcolor"
	"github.com/charmbracelet/lipgloss"
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

type ConstantRenderer interface{ Render(*RenderContext) string }
type ResponsiveRenderer interface {
	Render(*RenderContext) ConstantRenderer
}

type Defaulter interface {
	DefaultWidth()
	DefaultHeight()
}

// RenderContext represents current level in Renderer context tree (reactea's Component tree)
// In TUI terms you could say RenderContext is block
type RenderContext struct {
	display display

	width, height                 LengthNumber
	widthAttr, heightAttr         LengthAttribute
	occupiedWidth, occupiedHeight byte

	value string

	childrenContexts []*RenderContext

	parent *RenderContext
}

// Returns default SizeContext
// Zero SizeContext != Default SizeContext
func New() RenderContext {
	return RenderContext{
		// width:  FitContent,
		// height: FitContent,
	}
}

func ReverseSlice[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func Height[T Length](height T) {
	// switch any(height).(type) {
	// case LengthNumber:
	// 	height
	// }
}

func (rc *RenderContext) Propagate(child *RenderContext) {
	if rc.parent == nil {
		return
	}

	rc.parent.Propagate(rc)
}

func (rc *RenderContext) Add(renderer ConstantRenderer) *RenderContext {
	child := rc.Child()

	renderer.Render(child)

	rc.childrenContexts = append(rc.childrenContexts, child)

	return child
}

func (rc RenderContext) String() (result string) {
	if len(rc.childrenContexts) == 0 {
		return rc.value
	}

	elements := make(RenderedElements, 0, len(rc.childrenContexts))

	for _, child := range rc.childrenContexts {
		element := child.String()

		if debug {
			element = lipgloss.NewStyle().ColorWhitespace(true).Background(lipgloss.Color(randomcolor.GetRandomColorInHex())).Render(element)
		}

		switch child.display {
		case Block:

			elements = append(elements, RenderedElement{false, element})
		case Inline:
			elements = append(elements, RenderedElement{true, element})
		case InlineBlock:
		default:
			panic("unsupported display")
		}
	}

	return elements.Join()

	// if rc.display == Block {
	// 	widthPerElement := int(rc.width)
	// 	heightPerElement := int(rc.height) / len(elements)

	// 	fmt.Println(heightPerElement)
	// 	fmt.Println(widthPerElement)

	// 	for i := range elements {
	// 		elements[i] = lipgloss.NewStyle().Width(widthPerElement).Height(heightPerElement).Render(elements[i])
	// 		if debug {
	// 			fmt.Println(randomcolor.GetRandomColorInHex())
	// 			elements[i] = lipgloss.NewStyle().ColorWhitespace(true).Background(lipgloss.Color(randomcolor.GetRandomColorInHex())).Render(elements[i])
	// 		}
	// 	}

	// 	result = lipgloss.JoinVertical(lipgloss.Left, elements...)
	// } else {

	// }
	// return
}

func (rc *RenderContext) Child() *RenderContext {
	return &RenderContext{parent: rc}
}

func (rc *RenderContext) Value(value string) {
	rc.value = value
}

func (rc RenderContext) TreeString(indent string) string {
	result := indent + "- " + rc.display.String() + "\n"

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
