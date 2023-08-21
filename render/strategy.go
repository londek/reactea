package render

type Strategy uint8

const (
	// Cut extra content
	Cut Strategy = iota

	// Keep extra content
	Keep

	// Wrap extra content
	Wrap
)
