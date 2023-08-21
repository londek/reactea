package ruler

// display describes how element
// - is displayed with other siblings
// - is displaying children
type display int

func (d display) String() string {
	switch d {
	case Block:
		return "block"
	case Inline:
		return "inline"
	case InlineBlock:
		return "inlineblock"
	case Flex:
		return "flex"
	default:
		return "<unknown>"
	}
}

const (
	Inline display = iota
	Block
	InlineBlock
	Flex
)
