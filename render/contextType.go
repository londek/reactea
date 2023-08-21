package ruler

type Kind int

const (
	ContainerKind Kind = iota
	RenderableKind
)

func (kind Kind) String() string {
	switch kind {
	case ContainerKind:
		return "ContainerKind"
	case RenderableKind:
		return "RenderableKind"
	default:
		return "<unknown>"
	}
}
