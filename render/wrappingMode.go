package render

type WrappingMode uint8

const (
	// Cut extra content
	ClipContent WrappingMode = iota

	// Keep extra content
	KeepContent

	// Wrap extra content
	WrapContent
)
