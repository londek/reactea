package ruler

// display describes how element
// - is displayed with other siblings
// - is displaying children
type display int

func (d display) String() (result string) {
	switch d {
	case Block:
		result = "block"
	case Inline:
		result = "inline"
	case InlineBlock:
		result = "inlineblock"
	case Flex:
		result = "flex"
	}

	return
}

const (
	Inline display = iota
	Block
	InlineBlock
	Flex
)
