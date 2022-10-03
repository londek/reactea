package ruler

type Direction uint8

const (
	// Top to bottom / Left to right
	MinToMax Direction = iota

	// Bottom to top / Right to left
	MaxToMin
)
