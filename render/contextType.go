package ruler

type Kind int

const (
	ContainerKind Kind = iota
	RenderableKind
)

func (k Kind) String() string {
	switch k {
	case ContainerKind:
		return "ContainerKind"
	case RenderableKind:
		return "RenderableKind"
	default:
		return "<unknown>"
	}
}
