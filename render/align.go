package render

type Align uint8

const (
	// Min means align to minimum index:
	// - vertically: to the top
	// - horizontally: to the left
	Min Align = iota

	// Middle means align to middle
	Middle

	// Max means align to maximum index:
	// - vertically: to the bottom
	// - horizontally: to the right
	Max
)
