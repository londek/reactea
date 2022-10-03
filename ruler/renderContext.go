package ruler

import (
	"fmt"

	"github.com/AvraamMavridis/randomcolor"
	"github.com/charmbracelet/lipgloss"
)

var debug bool

func Debug() {
	debug = true
}

type Renderer = func(*RenderContext)

// RenderContext represents current level in Renderer context tree (reactea's Component tree)
// In TUI terms you could say RenderContext is block
type RenderContext struct {
	axis      Axis
	direction Direction

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

type SizeContextOption func(*RenderContext)

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

func (rc *RenderContext) Add(renderer Renderer, opts ...SizeContextOption) *RenderContext {
	child := rc.Child()

	rc.Apply(opts...)

	renderer(child)

	rc.childrenContexts = append(rc.childrenContexts, child)

	return child
}

func (rc *RenderContext) AddParagraph(text string, opts ...SizeContextOption) *RenderContext {
	return rc.Add(func(rc *RenderContext) {
		rc.Value(text)
	}, opts...)
}

func (rc RenderContext) String() (result string) {
	var elements []string

	if len(rc.childrenContexts) == 0 {
		return rc.value
	} else {
		for _, x := range rc.childrenContexts {
			elements = append(elements, x.String())
		}
	}

	if rc.direction == MaxToMin {
		ReverseSlice(elements)
	}

	fmt.Println(elements)

	if rc.axis == Vertical {
		widthPerElement := int(rc.width)
		heightPerElement := int(rc.height) / len(elements)

		fmt.Println(heightPerElement)
		fmt.Println(widthPerElement)

		for i := range elements {
			elements[i] = lipgloss.NewStyle().Width(widthPerElement).Height(heightPerElement).Render(elements[i])
			if debug {
				fmt.Println(randomcolor.GetRandomColorInHex())
				elements[i] = lipgloss.NewStyle().ColorWhitespace(true).Background(lipgloss.Color(randomcolor.GetRandomColorInHex())).Render(elements[i])
			}
		}

		result = lipgloss.JoinVertical(lipgloss.Left, elements...)
	} else {
		widthPerElement := int(rc.width) / len(elements)
		heightPerElement := int(rc.height)
		fmt.Println(heightPerElement)
		fmt.Println(widthPerElement)
		for i := range elements {
			elements[i] = lipgloss.NewStyle().Width(widthPerElement).Height(heightPerElement).Render(elements[i])
			if debug {
				fmt.Println(randomcolor.GetRandomColorInHex())
				elements[i] = lipgloss.NewStyle().ColorWhitespace(true).Background(lipgloss.Color(randomcolor.GetRandomColorInHex())).Render(elements[i])
			}
		}

		result = lipgloss.JoinHorizontal(lipgloss.Left, elements...)
	}

	return
}

func (rc *RenderContext) Child() *RenderContext {
	return &RenderContext{parent: rc}
}

func (rc *RenderContext) Value(value string) {
	rc.value = value
}

func (rc *RenderContext) Apply(opts ...SizeContextOption) {
	for _, opt := range opts {
		opt(rc)
	}
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
